/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Modifications made
 1. Removed code to handle many types of k8s object such as deployments,
		 pods etc.
 2. Remove namespace from Event.
 3. Modified #processItem to cast all interfaces to secrets and releases
    and to handle various types of Helm events.
*/

package controller

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/larderdev/kubewise/config"
	"github.com/larderdev/kubewise/driver"
	"github.com/larderdev/kubewise/handlers"
	"github.com/larderdev/kubewise/kwrelease"
	"github.com/larderdev/kubewise/utils"

	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const maxRetries = 5

var serverStartTime time.Time

type Event struct {
	key        string
	eventType  string
	secretType api_v1.SecretType
}

type Controller struct {
	clientset    kubernetes.Interface
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	eventHandler handlers.Handler
}

func Start(conf *config.Config, eventHandler handlers.Handler) {
	kubeClient := utils.GetClient()

	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return kubeClient.CoreV1().Secrets(conf.Namespace).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return kubeClient.CoreV1().Secrets(conf.Namespace).Watch(options)
			},
		},
		&api_v1.Secret{},
		0,
		cache.Indexers{},
	)

	c := newResourceController(kubeClient, eventHandler, informer)
	stopCh := make(chan struct{})
	defer close(stopCh)

	go c.Run(stopCh)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

func newResourceController(client kubernetes.Interface, eventHandler handlers.Handler, informer cache.SharedIndexInformer) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var newEvent Event
	var err error

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(secret interface{}) {
			newEvent.key, err = cache.MetaNamespaceKeyFunc(secret)
			newEvent.eventType = "create"
			newEvent.secretType = secret.(*api_v1.Secret).Type

			if err == nil {
				queue.Add(newEvent)
			}
		},

		UpdateFunc: func(secret, new interface{}) {
			newEvent.key, err = cache.MetaNamespaceKeyFunc(secret)
			newEvent.eventType = "update"
			newEvent.secretType = secret.(*api_v1.Secret).Type

			if err == nil {
				queue.Add(newEvent)
			}
		},

		DeleteFunc: func(secret interface{}) {
			newEvent.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(secret)
			newEvent.eventType = "delete"
			newEvent.secretType = secret.(*api_v1.Secret).Type

			if err == nil {
				queue.Add(newEvent)
			}
		},
	})

	return &Controller{
		clientset:    client,
		informer:     informer,
		queue:        queue,
		eventHandler: eventHandler,
	}
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	log.Println("Starting KubeWise controller")
	serverStartTime = time.Now().Local()

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	log.Println("KubeWise controller ready")

	wait.Until(c.runWorker, time.Second, stopCh)
}

func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

func (c *Controller) LastSyncResourceVersion() string {
	return c.informer.LastSyncResourceVersion()
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	newEvent, quit := c.queue.Get()

	if quit {
		return false
	}
	defer c.queue.Done(newEvent)

	err := c.processItem(newEvent.(Event))

	if err == nil {
		c.queue.Forget(newEvent)
	} else if c.queue.NumRequeues(newEvent) < maxRetries {
		log.Printf("Error processing %s (will retry): %v", newEvent.(Event).key, err)
		c.queue.AddRateLimited(newEvent)
	} else {
		log.Printf("Error processing %s (giving up): %v", newEvent.(Event).key, err)
		c.queue.Forget(newEvent)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) processItem(newEvent Event) error {
	object, _, err := c.informer.GetIndexer().GetByKey(newEvent.key)

	// GetByKey returns a nil object in the case where a Helm secret has been deleted. This means
	// we don't have access to the original secret at this point and can't inform the user about
	// which application has been successfully deleted.
	//
	// One approach to investigate in future would be to put the relevant details on the Event
	// and use that to provide the user information about what has been uninstalled.

	if err != nil {
		log.Fatalf("Error fetching secret with key %s from store: %v", newEvent.key, err)
		return err
	}

	// Uninstalling a Helm chart triggers a processItem but the secret has been deleted.
	// Without a nil check, we can see a panic when we type check the secret below.
	if object == nil {
		log.Println("Skipping nil secret", newEvent.eventType, "event for secret type:", newEvent.secretType)
		return nil
	}

	secret, ok := object.(*api_v1.Secret)

	if !ok {
		log.Println("Unable to cast 'object' (interface) as secret in", newEvent.eventType, "event:", object)
	}

	if secret.Type != "helm.sh/release.v1" {
		log.Println("Skipping non-helm secret", newEvent.eventType, "event:", secret.Type)
		return nil
	}

	currentRelease, err := driver.DecodeRelease(string(secret.Data["release"]))

	if err != nil {
		log.Fatalln("Error getting releaseData from secret", secret)
	}

	// This can be nil if, for example, this is the first time we are installing this chart.
	previousRelease := kwrelease.GetPreviousRelease(secret)

	switch newEvent.eventType {
	case "create":
		if secret.ObjectMeta.CreationTimestamp.Sub(serverStartTime).Seconds() > 0 {
			c.eventHandler.ObjectCreated(currentRelease, previousRelease)
			return nil
		}

	case "update":
		c.eventHandler.ObjectUpdated(currentRelease, previousRelease)
		return nil

	case "delete":
		c.eventHandler.ObjectDeleted(currentRelease, previousRelease)
		return nil
	}

	return nil
}

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
 1. Deleted superfluous code for getting object MetaData.
 2. Added GetClient function and made others private.
*/

package utils

import (
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// Pull in all the auth libs to ensure compatibility with many cluster providers.
	// https://github.com/kubernetes/client-go/issues/242#issuecomment-314642965
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// GetClient will read the kubectl from within or outside a cluster.
func GetClient() kubernetes.Interface {
	_, err := rest.InClusterConfig()
	if err != nil {
		return getClientOutOfCluster()
	}
	return getClientInsideCluster()
}

func getClientInsideCluster() kubernetes.Interface {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Can not get kubernetes config.")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Can not create kubernetes client.")
	}

	return clientset
}

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

func getClientOutOfCluster() kubernetes.Interface {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		log.Fatalf("Can not get kubernetes config.")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Can not get kubernetes config.")
	}

	return clientset
}

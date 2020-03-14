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

func GetClient() kubernetes.Interface {
	_, err := rest.InClusterConfig()
	if err != nil {
		return GetClientOutOfCluster()
	}
	return GetClientInsideCluster()
}

func GetClientInsideCluster() kubernetes.Interface {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Can not create kubernetes client: %v", err)
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

func GetClientOutOfCluster() kubernetes.Interface {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		log.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Can not get kubernetes config: %v", err)
	}

	return clientset
}

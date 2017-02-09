package main

import (
	"flag"
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/kubernetes/typed/batch/v2alpha1"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "/home/samuel/.kube/config", "absolute path to the kubeconfig file")
)

func main() {
	flag.Parse()
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	check(err)
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	check(err)
	for {
		pods, err := clientset.Core().Pods("").List(v1.ListOptions{})
		check(err)
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		for i, pod := range pods.Items {
			fmt.Printf("Pod %d: %s\n", i, pod.Name)
		}
		nss := clientset.Core().Namespaces()
		fmt.Printf("NS: %s\n", nss)
		time.Sleep(2 * time.Second)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

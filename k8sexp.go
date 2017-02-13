package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "/home/samuel/.kube/config", "absolute path to the kubeconfig file")
)

func main() {
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	check(err)
	clientset, err := kubernetes.NewForConfig(config)
	check(err)
	jobsList, err := clientset.BatchV1Client.Jobs("default").List(v1.ListOptions{})
	for i, job := range jobsList.Items {
		fmt.Printf("Job %d: %s\n", i, job.Name)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

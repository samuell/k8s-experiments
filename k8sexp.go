package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
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

	jobs := clientset.BatchV1Client.Jobs("default")
	fmt.Printf("Jobs: %v\n", jobs)

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

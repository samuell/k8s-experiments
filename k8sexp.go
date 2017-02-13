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

	// Access jobs. We can't do it all in one line, since we need to receive the
	// errors and manage thgem appropriately
	batchClient := clientset.BatchV1Client
	jobsClient := batchClient.Jobs("default")
	piJob, err := jobsClient.Get("pi")
	check(err)
	fmt.Printf("piJob Name: %v\n", piJob.Name)

	jobsList, err := jobsClient.List(v1.ListOptions{})
	check(err)

	// Loop over all jobs and print their name
	for i, job := range jobsList.Items {
		fmt.Printf("Job %d: %s\n", i, job.Name)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

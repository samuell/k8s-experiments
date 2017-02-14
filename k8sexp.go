package main

import (
	"flag"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/unversioned"
	k8sapi "k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "/home/samuel/.kube/config", "absolute path to the kubeconfig file")
	trueVal    = true
	falseVal   = false
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

	jobsList, err := jobsClient.List(k8sapi.ListOptions{})
	check(err)

	// Loop over all jobs and print their name
	for i, job := range jobsList.Items {
		fmt.Printf("Job %d: %s\n", i, job.Name)
	}

	// For an example of how to create jobs, see this file:
	// https://github.com/pachyderm/pachyderm/blob/805e63e561a9eab4a9c52216f228f0f421714f3b/src/server/pps/server/api_server.go#L2320-L2345
	batchJob := &batchv1.Job{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Job",
			APIVersion: "v1",
		},
		ObjectMeta: k8sapi.ObjectMeta{
			Name:   "k8sexp-testjob",
			Labels: make(map[string]string),
		},
		Spec: batchv1.JobSpec{
			// Optional: Parallelism:,
			// Optional: Completions:,
			// Optional: ActiveDeadlineSeconds:,
			// Optional: Selector:,
			// Optional: ManualSelector:,
			Template: k8sapi.PodTemplateSpec{
				ObjectMeta: k8sapi.ObjectMeta{
					Name:   "k8sexp-testpod",
					Labels: make(map[string]string),
				},
				Spec: k8sapi.PodSpec{
					InitContainers: []k8sapi.Container{}, // Doesn't seem obligatory(?)...
					Containers: []k8sapi.Container{
						{
							Name:    "k8sexp-testimg",
							Image:   "perl",
							Command: []string{"sleep", "10"},
							SecurityContext: &k8sapi.SecurityContext{
								Privileged: &falseVal,
							},
							ImagePullPolicy: k8sapi.PullPolicy(k8sapi.PullIfNotPresent),
							Env:             []k8sapi.EnvVar{},
							VolumeMounts:    []k8sapi.VolumeMount{},
						},
					},
					RestartPolicy:    k8sapi.RestartPolicyOnFailure,
					Volumes:          []k8sapi.Volume{},
					ImagePullSecrets: []k8sapi.LocalObjectReference{},
				},
			},
		},
		// Optional, not used by pach: JobStatus:,
	}

	newJob, err := jobsClient.Create(batchJob)
	check(err)

	fmt.Println("New job name: ", newJob.Name)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

package main

import (
	"flag"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
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
	// From https://github.com/pachyderm/pachyderm/blob/805e63e561a9eab4a9c52216f228f0f421714f3b/src/server/pps/server/api_server.go#L2320-L2345
	//
	// import (
	//		"k8s.io/kubernetes/pkg/api"
	//		"k8s.io/kubernetes/pkg/api/unversioned"
	//		"k8s.io/kubernetes/pkg/apis/batch"
	//		kube "k8s.io/kubernetes/pkg/client/unversioned"
	//		kube_labels "k8s.io/kubernetes/pkg/labels"
	// )
	//
	// jobStructure := &batch.Job{
	// 	TypeMeta: unversioned.TypeMeta{
	// 		Kind:       "Job",
	// 		APIVersion: "v1",
	// 	},
	// 	ObjectMeta: api.ObjectMeta{
	// 		Name:   jobInfo.JobID,
	// 		Labels: options.labels,
	// 	},
	// 	Spec: batch.JobSpec{
	// 		ManualSelector: &trueVal,
	// 		Selector: &unversioned.LabelSelector{
	// 			MatchLabels: options.labels,
	// 		},
	// 		Parallelism: &options.parallelism,
	// 		Completions: &options.parallelism,
	// 		Template: api.PodTemplateSpec{
	// 			ObjectMeta: api.ObjectMeta{
	// 				Name:   jobInfo.JobID,
	// 				Labels: options.labels,
	// 			},
	// 			Spec: podSpec(options, jobInfo.JobID, "Never"),
	// 		},
	// 	},
	// }
	//
	// func podSpec(options *jobOptions, jobID string, restartPolicy api.RestartPolicy) api.PodSpec {
	// 	return api.PodSpec{
	// 		InitContainers: []api.Container{
	// 			{
	// 				Name:            "init",
	// 				Image:           options.jobShimImage,
	// 				Command:         []string{"/pach/job-shim.sh"},
	// 				ImagePullPolicy: api.PullPolicy(options.jobImagePullPolicy),
	// 				Env:             options.jobEnv,
	// 				VolumeMounts:    options.volumeMounts,
	// 			},
	// 		},
	// 		Containers: []api.Container{
	// 			{
	// 				Name:    "user",
	// 				Image:   options.userImage,
	// 				Command: []string{"/pach-bin/guest.sh", jobID},
	// 				SecurityContext: &api.SecurityContext{
	// 					Privileged: &trueVal, // god is this dumb
	// 				},
	// 				ImagePullPolicy: api.PullPolicy(options.jobImagePullPolicy),
	// 				Env:             options.jobEnv,
	// 				VolumeMounts:    options.volumeMounts,
	// 			},
	// 		},
	// 		RestartPolicy:    restartPolicy,
	// 		Volumes:          options.volumes,
	// 		ImagePullSecrets: options.imagePullSecrets,
	// 	}
	// }

	batchJob := &batchv1.Job{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Job",
			APIVersion: "v1",
		},
		ObjectMeta: api.ObjectMeta{
			Name:   jobInfo.JobID, // TODO: What is jobInfo?
			Labels: options.labels,
		},
		Spec: batchv1.JobSpec{
			// Optional: Parallelism:,
			// Optional: Completions:,
			// Optional: ActiveDeadlineSeconds:,
			// Optional: Selector:,
			// Optional: ManualSelector:,
			Template: api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Name:   jobInfo.JobID,
					Labels: options.labels,
				},
				Spec: api.PodSpec{
					InitContainers: []api.Container{
						{
							Name:            "init",
							Image:           options.jobShimImage,
							Command:         []string{"/pach/job-shim.sh"},
							ImagePullPolicy: api.PullPolicy(options.jobImagePullPolicy),
							Env:             options.jobEnv,
							VolumeMounts:    options.volumeMounts,
						},
					},
					Containers: []api.Container{
						{
							Name:    "user",
							Image:   options.userImage,
							Command: []string{"/pach-bin/guest.sh", jobID},
							SecurityContext: &api.SecurityContext{
								Privileged: &trueVal, // god is this dumb
							},
							ImagePullPolicy: api.PullPolicy(options.jobImagePullPolicy),
							Env:             options.jobEnv,
							VolumeMounts:    options.volumeMounts,
						},
					},
					RestartPolicy:    restartPolicy,
					Volumes:          options.volumes,
					ImagePullSecrets: options.imagePullSecrets,
				},
			},
		},
		// Optional, not used by pachy: JobStatus:,
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

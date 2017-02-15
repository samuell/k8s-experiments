package main

import (
	"flag"
	"fmt"

	"k8s.io/client-go/kubernetes"
	apiUnver "k8s.io/client-go/pkg/api/unversioned"
	api "k8s.io/client-go/pkg/api/v1"
	batchapi "k8s.io/client-go/pkg/apis/batch/v1"
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

	jobsList, err := jobsClient.List(api.ListOptions{})
	check(err)

	// Loop over all jobs and print their name
	for i, job := range jobsList.Items {
		fmt.Printf("Job %d: %s\n", i, job.Name)
	}

	//storageQuantity1Gi, err := resource.ParseQuantity("1Gi")
	//check(err)

	//	k8sexpVolume := &api.PersistentVolume{
	//		ObjectMeta: api.ObjectMeta{
	//			Name: "k8sexp-volume",
	//		},
	//		Spec: api.PersistentVolumeSpec{
	//			Capacity: api.ResourceList{
	//				api.ResourceStorage: storageQuantity1Gi,
	//			},
	//			AccessModes: []api.PersistentVolumeAccessMode{
	//				api.PersistentVolumeAccessMode("ReadWriteMany"),
	//			},
	//			PersistentVolumeReclaimPolicy: api.PersistentVolumeReclaimRecycle,
	//		},
	//	}
	//
	//	k8sexpVolumeClaim := &api.PersistentVolumeClaim{}
	//
	//	fmt.Printf("Volume: %v\n\nVolumeClaim: %v\n", k8sexpVolume, k8sexpVolumeClaim)

	// For an example of how to create jobs, see this file:
	// https://github.com/pachyderm/pachyderm/blob/805e63/src/server/pps/server/api_server.go#L2320-L2345
	batchJob := &batchapi.Job{
		TypeMeta: apiUnver.TypeMeta{
			Kind:       "Job",
			APIVersion: "v1",
		},
		ObjectMeta: api.ObjectMeta{
			Name:   "k8sexp-testjob",
			Labels: make(map[string]string),
		},
		Spec: batchapi.JobSpec{
			// Optional: Parallelism:,
			// Optional: Completions:,
			// Optional: ActiveDeadlineSeconds:,
			// Optional: Selector:,
			// Optional: ManualSelector:,
			Template: api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Name:   "k8sexp-testpod",
					Labels: make(map[string]string),
				},
				Spec: api.PodSpec{
					InitContainers: []api.Container{}, // Doesn't seem obligatory(?)...
					Containers: []api.Container{
						{
							Name:    "k8sexp-testimg",
							Image:   "perl",
							Command: []string{"sh", "-c", "echo hej > /k8sexp-data/hej.txt"},
							SecurityContext: &api.SecurityContext{
								Privileged: &falseVal,
							},
							ImagePullPolicy: api.PullPolicy(api.PullIfNotPresent),
							Env:             []api.EnvVar{},
							VolumeMounts: []api.VolumeMount{
								api.VolumeMount{
									Name:      "k8sexp-testvol",
									MountPath: "/k8sexp-data",
								},
							},
						},
					},
					RestartPolicy:    api.RestartPolicyOnFailure,
					ImagePullSecrets: []api.LocalObjectReference{},
					Volumes: []api.Volume{
						api.Volume{
							Name: "k8sexp-testvol",
							VolumeSource: api.VolumeSource{
								HostPath: &api.HostPathVolumeSource{
									Path: "/data",
								},
							},
						},
					},
				},
			},
		},
		// Optional, not used by pach: JobStatus:,
	}
	//Volumes: []Volumes{
	//	api.PersistentVolume{
	//		ObjectMeta: api.ObjectMeta{
	//			Name: "k8sexp-testvol",
	//		},
	//		Spec: api.PersistentVolumeSpec{
	//			AccessModes: api.ReadWriteOnce,
	//		},
	//	},
	//},

	newJob, err := jobsClient.Create(batchJob)
	check(err)

	fmt.Println("New job name: ", newJob.Name)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

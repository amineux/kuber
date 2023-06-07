
package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Path to the kubeconfig file
	kubeconfig := flag.String("kubeconfig", filepath.Join(
		homeDir(), ".kube", "config"), "Path to the kubeconfig file")
	flag.Parse()

	// Build the client configuration
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// Create the Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// List all Pods in the default namespace
	pods, err := clientset.CoreV1().Pods("default").List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Pods in the default namespace:")
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}

	// Create a new Pod
	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example-pod",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "example-container",
					Image: "nginx",
				},
			},
		},
	}

	createdPod, err := clientset.CoreV1().Pods("default").Create(newPod)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created Pod: %s\n", createdPod.Name)
}

// Helper function to get the user's home directory
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // Windows
}

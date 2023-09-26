package controller

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func getAll() {

	namespace := "myns2"
	myConfig := "/home/shi/.kube/config"
	// Load the kubeconfig file and create a Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", myConfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		//os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		//os.Exit(1)
	}
	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing services: %v\n", err)
		//os.Exit(1)
	}

	// Print the names of the services.
	fmt.Printf("Services in namespace %s:\n", namespace)
	for _, service := range services.Items {
		fmt.Printf("- %s\n", service.Name)
	}
}

func getPod(podName string, namespace string) error {
	// Define the service name and namespace you want to check
	//podName := "dbsvc"
	//namespace := "myns2"

	myConfig := "/home/shi/.kube/config"
	// Load the kubeconfig file and create a Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", myConfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		//os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		//os.Exit(1)
	}

	myDep, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil || myDep == nil {
		fmt.Printf("Error fetching statefulset: %v\n", err)
		//os.Exit(1)
		return err
	}

	return nil
}
func searchSvc() (v1.Endpoints, error) {
	// Define the service name and namespace you want to check
	serviceName := "dbsvc"
	namespace := "myns2" // Replace with the appropriate namespace

	// Set up the kubeconfig file path (or use in-cluster configuration)
	//var kubeconfig string
	//home := homeDir()
	//flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")
	//flag.Parse()
	myConfig := "/home/shi/.kube/config"
	// Load the kubeconfig file and create a Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", myConfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		//os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		//os.Exit(1)
	}

	// Check the availability of the service
	//endpoint, err := clientset.CoreV1().Endpoints(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	endpoint, err := clientset.CoreV1().Endpoints(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})

	if err != nil {
		fmt.Printf("Error fetching service endpoints: %v\n", err)
		//os.Exit(1)
		return *endpoint, err
	}
	//if len(endpoint.Subsets) > 0 {
	//	fmt.Printf("Service %s is available in namespace %s\n", serviceName, namespace)
	//} else {
	//	fmt.Printf("Service %s is not available in namespace %s\n", serviceName, namespace)
	//}
	return *endpoint, nil
}

// homeDir returns the user's home directory.
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // Windows
}

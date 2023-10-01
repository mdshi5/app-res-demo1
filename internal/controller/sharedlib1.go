package controller

import (
	"context"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
)

func (r *EntrycrReconciler) depInfo(deploymentName string, namespace string, l logr.Logger) (bool, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		if home := homedir.HomeDir(); home != "" {
			kubeconfigPath := filepath.Join(home, ".kube", "config")
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
			if err != nil {
				l.Info("Error loading kubeconfig")
				return false, err
			}
		} else {
			l.Info("Error loading kubeconfig")
			return false, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		l.Info("Error creating Kubernetes client")
		return false, err
	}

	deploymentsClient := clientset.AppsV1().StatefulSets(namespace)

	// Define a retry strategy with exponential backoff
	retryStrategy := wait.Backoff{
		Steps:    5,
		Duration: 5 * time.Second,
		Factor:   2.0,
		Jitter:   0.1,
	}

	// Retry checking readiness until it succeeds or the retry strategy gives up
	err = wait.ExponentialBackoff(retryStrategy, func() (bool, error) {
		deployment, err := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if deployment.Status.ReadyReplicas == deployment.Status.Replicas {
			l.Info("Deployment is ready.")
			return true, nil
		} else {
			l.Info("Deployment is not yet ready.")
		}
		return false, nil

	})
	if err != nil {
		l.Info("Error checking readiness")
		return false, err
	} else {
		return true, nil
	}

}

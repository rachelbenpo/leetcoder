package services

import (
	"context"
	"fmt"
	"io"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// takes a docker image and runs it inside a kubernetes pod. returns the pod output
func runImageInsideK8S(imageName string) (string, error) {

	// get k8s client
	client, err := getClient()
	if err != nil {
		return "", err
	}

	// Run Kubernetes pod
	podName, err := createPod(imageName, client)
	if err != nil {
		fmt.Println("Error creating Kubernetes pod:", err)
		return "", err
	}
	fmt.Println("ran Kubernetes pod: ", podName)

	// wait for pod to complete
	waitForPodCompletion(podName, client)

	// Get pod output
	output, err := getPodOutput(podName, client)
	if err != nil {
		fmt.Println("Error getting pod output:", err)
		return "", err
	}

	// Remove the pod
	err = removePod(podName, client)
	if err != nil {
		fmt.Println("Error cleaning up:", err)
		return output, err
	}

	return output, nil
}

// create a Kubernetes pod using the specified image name. returns pod name
func createPod(imageName string, clientset *kubernetes.Clientset) (string, error) {

	// set pod configuration
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "checking-pod-",
		},
		Spec: v1.PodSpec{
			RestartPolicy: v1.RestartPolicyNever,
			Containers: []v1.Container{
				{
					Name:  "checking-container",
					Image: imageName,
				},
			},
		},
	}

	// create the pod
	pod, err := clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	return pod.Name, err
}

// wait until the pod is completed
func waitForPodCompletion(podName string, clientset *kubernetes.Clientset) error {

	pollingInterval := 2 * time.Second
	maxWaitTimeout := 6 * time.Minute

	return wait.PollImmediate(pollingInterval, maxWaitTimeout, func() (done bool, err error) {
		pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		// Check if the container has terminated
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == "checking-container" && containerStatus.State.Terminated != nil {
				return true, nil
			}
		}
		return false, nil
	})
}

// get the output of a Kubernetes pod
func getPodOutput(podName string, clientset *kubernetes.Clientset) (string, error) {

	// Get Pod logs
	podLogs, err := clientset.CoreV1().Pods("default").GetLogs(podName, &v1.PodLogOptions{}).Stream(context.Background())
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	// Read Pod logs
	output, err := io.ReadAll(podLogs)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// delete pod
func removePod(podName string, clientset *kubernetes.Clientset) error {

	err := clientset.CoreV1().Pods("default").Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

// create a Kubernetes client
func getClient() (*kubernetes.Clientset, error) {

	// get config path
	home := homedir.HomeDir()
	kubeConfigPath := home + "/.kube/config"

	// Load kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

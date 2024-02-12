package services

import (
	"context"
	"fmt"
	"io"
	"time"
	// "path/filepath"

	// batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func manageK8s(dockerCode string) (string, error) {

	imageName := "checking-container"
	imageName = "ghcr.io/rachelbenpo/check-code-try:2-false"

	// // Build Docker image
	// err := buildImage(dockerCode, imageName)
	// if err != nil {
	// fmt.Println("Error building Docker image:", err)
	// return "", err
	// }

	fmt.Println("built Docker image: ", imageName)

	// Run Kubernetes pod
	podName, err := runPod(imageName)
	if err != nil {
		fmt.Println("Error running Kubernetes pod:", err)
		return "", err
	}

	fmt.Println("ran Kubernetes pod: ", podName)

	// Get pod output
	output, err := getPodOutput(podName)
	if err != nil {
		fmt.Println("Error getting pod output:", err)
		return "", err
	}

	fmt.Println("got pod output: ", output)

	// Remove the pod and image
	err = removePodAndImage(podName, imageName)
	if err != nil {
		fmt.Println("Error cleaning up:", err)
		return "", err
	}

	fmt.Println("cleaned up")

	return output, nil
}

func buildPod(imageName string) error {

	// Load kubeconfig
	home := homedir.HomeDir()
	kubeconfig := fmt.Sprintf("%s/.kube/config", home)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// Create Pod
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod3-" + imageName,
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-" + imageName,
					Image: imageName,
				},
			},
		},
	}

	_, err = clientset.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Println("Kubernetes Pod built successfully:", pod.Name)

	return nil
}

// run a Kubernetes pod using the specified image name. returns pod name
func runPod(imageName string) (string, error) {

	// Load kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", getKubeConfigPath())
	if err != nil {
		return "", err
	}
	fmt.Print("loaded kubeconfig\n")

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}
	fmt.Print("created clientset\n")

	// Create Pod
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

	pod, err = clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	fmt.Print("created pod\n")

	podName := pod.Name

	pollingInterval := 2 * time.Second
	maxWaitTimeout := 30 * time.Second

	// Wait for the container to terminate
	err = wait.PollImmediate(pollingInterval, maxWaitTimeout, func() (done bool, err error) {
		pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		fmt.Printf("waiting for pod\n")

		// Check if the container has terminated
		for _, containerStatus := range pod.Status.ContainerStatuses {
			fmt.Println(containerStatus.Name)
			fmt.Println(containerStatus.State)

			if containerStatus.Name == "checking-container" {
				if containerStatus.State.Terminated != nil {
					return true, nil
				}
			}
		}

		return false, nil
	})

	// err = wait.PollImmediate(wait.ForeverTestTimeout, wait.ForeverTestTimeout, func() (done bool, err error) {
	// 	pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), podName, metav1.GetOptions{})
	// 	if err != nil {
	// 		return false, err
	// 	}

	// 	if pod.Status.Phase == v1.PodSucceeded || pod.Status.Phase == v1.PodFailed {
	// 		return true, nil
	// 	}

	// 	return false, nil
	// })

	return podName, err
}

// get the output of a Kubernetes pod
func getPodOutput(podName string) (string, error) {

	// Load kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", getKubeConfigPath())
	if err != nil {
		return "", err
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	// Get Pod logs
	podLogs, err := clientset.CoreV1().Pods("default").GetLogs(podName, &v1.PodLogOptions{}).Stream(context.TODO())
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	// Read Pod logs
	var outputBytes []byte
	buf := make([]byte, 1024)
	for {
		n, err := podLogs.Read(buf)
		if n == 0 && err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			return "", err
		}
		outputBytes = append(outputBytes, buf[:n]...)
	}

	return string(outputBytes), nil
}

// delete pod and image
func removePodAndImage(podName, imageName string) error {

	// Load kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", getKubeConfigPath())
	if err != nil {
		return err
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// Delete pod
	err = clientset.CoreV1().Pods("default").Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	// TODO: Delete image in K8s environment if needed

	return nil
}

// TODO: is this the correct path?
func getKubeConfigPath() string {
	home := homedir.HomeDir()
	return home + "/.kube/config"
}

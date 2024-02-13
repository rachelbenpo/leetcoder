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

	// get k8s client
	client, err := getClient()
	if err != nil {
		return "", err
	}

	imageName := "checking-container"

	// Build Docker image
	err = buildImage(dockerCode, imageName)
	if err != nil {
		fmt.Println("Error building Docker image:", err)
		return "", err
	}
	fmt.Println("built Docker image: ", imageName)

	// Push the Docker image to GitHub Container Registry
	_, err = pushImage(imageName)
	if err != nil {
		fmt.Println("Error pushing image to registry:", err)
		return "", err
	}
	fmt.Println("pushed image to registry: ", imageName)

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

	fmt.Println("got pod output: ", output)

	// Remove the pod and image
	// err = removePodAndImage(podName, imageUrl, client)
	// if err != nil {
	// 	fmt.Println("Error cleaning up:", err)
	// 	return "", err
	// }

	// fmt.Println("cleaned up")

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
	fmt.Print("created pod\n")

	return pod.Name, err
}

// TOREMOVE if not needed (same as createPod function)
func buildPod(imageName string, clientset *kubernetes.Clientset) error {

	// Create Pod
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-" + imageName,
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

	_, err := clientset.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Println("Kubernetes Pod built successfully:", pod.Name)

	return nil
}

// wait until the pod is completed
func waitForPodCompletion(podName string, clientset *kubernetes.Clientset) error {

	pollingInterval := 2 * time.Second
	maxWaitTimeout := 3 * time.Minute

	return wait.PollImmediate(pollingInterval, maxWaitTimeout, func() (done bool, err error) {
		pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		fmt.Printf("waiting for pod\n")

		// Check if the container has terminated
		for _, containerStatus := range pod.Status.ContainerStatuses {

			fmt.Println(containerStatus.Name)
			fmt.Println(containerStatus.State)

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

// TODO: Delete image in K8s environment if needed. now the function only removes the pod.
// delete pod and image
func removePodAndImage(podName, imageName string, clientset *kubernetes.Clientset) error {

	// Delete pod
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

// TOREMOVE----
// TODO: is this the correct path?
// func getKubeConfigPath() string {
// 	home := homedir.HomeDir()
// 	return home + "/.kube/config"
// }

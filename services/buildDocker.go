package services

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func manageDocker(dockerCode string) (string, error) {

	imageName := "checking-container"

	// Build Docker image
	err := buildImage(dockerCode, imageName)
	if err != nil {
		fmt.Println("Error building Docker image:", err)
		return "", err
	}

	// Run Docker container
	containerID, err := runContainer(imageName)
	if err != nil {
		fmt.Println("Error running Docker container:", err)
		return "", err
	}

	// Get container output
	output, err := getContainerOutput(containerID)
	if err != nil {
		fmt.Println("Error getting container output:", err)
		return "", err
	}

	// remove the container and image
	err = removeContainerAndImage(containerID, imageName)
	if err != nil {
		fmt.Println("Error cleaning up:", err)
		return "", err
	}

	fmt.Println("output:", output)

	return output, nil
}

// build a docker image from the dockerfile
func buildImage(dockerfileContent string, imageName string) error {

	// Create a Docker client
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Create a tar archive from the Dockerfile content
	tarball, err := archive.Generate("Dockerfile", dockerfileContent)
	if err != nil {
		return err
	}

	// Specify build options
	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
	}

	// Build Docker image
	buildResponse, err := cli.ImageBuild(ctx, tarball, buildOptions)
	if err != nil {
		return err
	}
	defer buildResponse.Body.Close()

	// Print build output
	_, err = io.Copy(os.Stdout, buildResponse.Body)
	if err != nil {
		return err
	}

	fmt.Println("Docker image built successfully:", imageName)

	return nil
}

// run a Docker container using the specified image name. returns container ID
func runContainer(imageName string) (string, error) {

	ctx := context.Background()

	// Create a Docker client
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return "", err
	}

	// Create a new container
	resp, err := cli.ContainerCreate(ctx, &container.Config{Image: imageName}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	// Return container ID
	return resp.ID, nil
}

// get the output of the container
func getContainerOutput(containerID string) (string, error) {

	ctx := context.Background()

	// Create a Docker client
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return "", err
	}

	// Get the container logs
	out, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Read the output from the logs
	output, err := io.ReadAll(out)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// delete image and container
func removeContainerAndImage(containerID, imageName string) error {

	ctx := context.Background()

	// Create a Docker client
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}

	// Remove container
	err = cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
	if err != nil {
		return err
	}

	// Remove image
	_, err = cli.ImageRemove(ctx, imageName, types.ImageRemoveOptions{Force: true})
	if err != nil {
		return err
	}

	return nil
}

// func startDockerEngine() error {
// 	// Check if Docker is running
// 	_, err := exec.Command("docker", "info").CombinedOutput()
// 	if err == nil {
// 		return nil
// 	}
// 	// Docker is not running, try to start it
// 	startCmd := exec.Command("service", "docker", "start")
// 	startCmdOutput, startCmdErr := startCmd.CombinedOutput()
// 	if startCmdErr != nil {
// 		return fmt.Errorf("error starting Docker engine: %v\n%s", startCmdErr, startCmdOutput)
// 	}
// }

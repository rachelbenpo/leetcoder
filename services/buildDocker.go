package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"

	"leetcoder/config"
)

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

// push a docker image to github container registry
func pushImage(imageName string) (string, error) {

	// Create a Docker client
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return "", err
	}

	// login to ghcr.io
	authConfig := registry.AuthConfig{
		Username: config.UserName,
		Password: config.Token,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	// Push the image to GitHub Container Registry
	resp, err := cli.ImagePush(context.Background(), imageName, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		fmt.Println("error pushing Docker image: ", imageName, " ", err)
		return "", err
	}
	defer resp.Close()

	// Print the response message
	body, err := io.ReadAll(resp)
	if err != nil {
		return "", err
	}
	fmt.Println("Image push response:", string(body))

	return imageName, nil
}

// TOREMOVE: from here until end of file
// since don't need the docker functionality, only k8s

// run the test code in a docker container
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

// get the output of a Docker container
func getContainerOutput(containerID string) (string, error) {

	// Run "docker logs" command and to get the container's output
	cmd := exec.Command("docker", "logs", containerID)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output)[0 : len(output)-1], nil
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

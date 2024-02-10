package services

import (
	"context"
	"fmt"
	"io"
	"os"
	//"io/ioutil"
	//"os/exec"

	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// TODO:
// func manageDockerJS(code string) (string, error) {


// not good enough - error building docker image
func manageDockerPython(code string) (string, error) {

	// Build Docker image with Python code
	imageName := "python-container"
	err := buildImage(code, imageName)
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

	return output, nil
}

func buildImage(pythonCode, imageName string) error {

	// Create temp directory to store the Dockerfile and Python code
	tempDir, err := os.MkdirTemp("", "docker-example")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// Write Dockerfile
	dockerfileContent := fmt.Sprintf(`
FROM python:3
WORKDIR /app
RUN echo '%s' > script.py
CMD ["python", "script.py"]
`, strings.ReplaceAll(pythonCode, "'", `'"'"'`))

	dockerfilePath := filepath.Join(tempDir, "Dockerfile")
	err = os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644)
	if err != nil {
		return err
	}

	fmt.Print("after writing dockerfile", dockerfileContent)

	// Build Docker image
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}

	fmt.Print("after build docker image: ", dockerfilePath)

	//buildContext := filepath.Dir(dockerfilePath)
	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
	}

	buildResponse, err := cli.ImageBuild(ctx, nil, buildOptions)
	if err != nil {
		return err
	}
	defer buildResponse.Body.Close()

	// Copy build output to stdout
	_, err = io.Copy(os.Stdout, buildResponse.Body)
	if err != nil {
		return err
	}

	return nil
}

// run the image of the user's code
func runContainer(imageName string) (string, error) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return "", err
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
		},
		nil,
		nil,
		nil,
		"",
	)
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

// get the output of the container - if the answer is correct or not
func getContainerOutput(containerID string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return "", err
	}

	out, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}
	defer out.Close()

	output, err := io.ReadAll(out)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// delete image and container
func removeContainerAndImage(containerID, imageName string) error {
	ctx := context.Background()
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

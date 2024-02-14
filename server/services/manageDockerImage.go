package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"

	"leetcode-server/config"
)

// create dockerfile content string for testing the answer
func createDockerfileContent(code, lang string) (string, error) {

	if lang == "python" {
		return createPythonDockerfileContent(code), nil
	}

	if lang == "javascript" || lang == "js" {
		return createJSDockerfileContent(code), nil
	}
	return "", fmt.Errorf("code language is not supported: ", lang)
}

// create dockerfile content string for running python code
func createPythonDockerfileContent(pythonCode string) string {

	dockerfileContent := fmt.Sprintf(`
FROM python:3
WORKDIR /app
RUN echo '%s' > script.py
CMD ["python", "script.py"]
`, strings.ReplaceAll(pythonCode, "'", `'"'"'`))

	return dockerfileContent
}

// create dockerfile content string for running JavaScript code
func createJSDockerfileContent(jsCode string) string {
	dockerfileContent := fmt.Sprintf(`
FROM node:14
WORKDIR /app
RUN echo '%s' > script.js
CMD ["node", "script.js"]
`, strings.ReplaceAll(jsCode, "'", `'"'"'`))

	return dockerfileContent
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

// delete image and container
func removeImage(imageName string) error {

	ctx := context.Background()

	// Create a Docker client
	cli, err := client.NewClientWithOpts()
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

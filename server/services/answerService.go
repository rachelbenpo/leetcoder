package services

import (
	"fmt"
	"leetcode-server/config"
	"leetcode-server/models"
)

// check if user's answer is correct or not, using docker and kubernetes
func CheckAnswer(ans models.Answer, q models.Question) (string, error) {

	imageName := "ghcr.io/" + config.UserName + "/checking-container"

	// build code for testing answer with all test cases
	testCode, err := buildTestCode(ans, q)
	if err != nil {
		return "", fmt.Errorf("error building test code", err)
	}

	// create dockerfile content
	dockerCode, err := createDockerfileContent(testCode, ans.Lang)
	if err != nil {
		return "", fmt.Errorf("error creating dockerfile content", err)
	}

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

	// run the image inside k8s
	answer, err := runImageInsideK8S(imageName)
	if err != nil {
		fmt.Print("error running test code in k8s", err)
		return "", err
	}

	// remove the image
	err = removeImage(imageName)
	if err!= nil {
        fmt.Print("error removing image", err)
        return answer[:len(answer)-1], err
    }

	fmt.Println("answer: ", answer)

	return answer[:len(answer)-1], nil
}

package services

import (
	"leetcoder/models"
)


// check if user's answer is correct or not
func CheckAnswer(ans models.Answer, q models.Question) (string, error) {

	// build code for testing answer with all test cases
	testCode, err := buildTestCode(ans, q)
	if err != nil {
		return "", err
	}

	// build dockerfile
	dockerCode, err := buildDockerfile(testCode, ans.Lang)
	if err != nil {
		return "", err
	}

	// run the test code inside docker container
	answer, err := manageDocker(dockerCode)
	if err != nil {
		return "", err
	}

	return answer, nil
}

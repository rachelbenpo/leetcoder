package services

import (
	"fmt"
	"leetcoder/models"
)

// check if user's answer is correct or not
func CheckAnswer(ans models.Answer, q models.Question) (string, error) {

	// build code for testing answer with all test cases
	testCode, err := buildTestCode(ans, q)
	if err != nil {
		fmt.Print("error building test code", err)
		return "", err
	}
	
	fmt.Println("testCode: ", testCode)

	// build dockerfile
	dockerCode, err := buildDockerfile(testCode, ans.Lang)
	if err != nil {
		fmt.Print("error building dockerfile", err)
		return "", err
	}

	fmt.Println("dockerCode: ", dockerCode)

	// run the test code inside k8s container
	answer, err := manageK8s(dockerCode)
	if err != nil {
		fmt.Print("error running test code in k8s", err)
		return "", err
	}
	
	fmt.Println("answer: ", answer)

	return answer[:len(answer)-1], nil
}

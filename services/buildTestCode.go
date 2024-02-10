package services

import (
	"fmt"
	"leetcoder/models"
)

func CheckAnswer(ans models.Answer, q models.Question) (string, error) {

	code, err := buildCheckCode(ans, q)
	if err != nil {
		return "", err
	}

	answer, err := manageDockerPython(code)
	if err != nil {
		return "", err
	}

	fmt.Printf(answer)

	return answer, nil
}

func buildCheckCode(ans models.Answer, q models.Question) (string, error) {

	if ans.Lang == "python" {
		return buildPythonCode(ans, q)
	}

	if ans.Lang == "javascript" || ans.Lang == "js" {
		return buildJsCode(ans, q)
	}
	return "", fmt.Errorf("code language is not supported")
}

// generates a python code that runs the user's code against all test cases
func buildPythonCode(ans models.Answer, q models.Question) (string, error) {

	codeToExec :=
		`import json
	` + ans.Code +
			`def main():
		inputs = [` + q.TestCases[0].Input

	for _, t := range q.TestCases[1:] {
		codeToExec += ` ,` + t.Input
	}

	codeToExec += `]
		outputs = [` + q.TestCases[0].Output

	for _, t := range q.TestCases[1:] {
		codeToExec += ` ,` + t.Output
	}

	codeToExec +=
		`]
    	for i in range(len(inputs)):
        	ans = func_user(inputs[i])
        	json1 = json.dumps(ans)
        	json2 = json.dumps(outputs[i])
        	if json1 != json2:
            	return False
    	return True
		
		
	if __name__ == "__main__":
	main()`

	return codeToExec, nil
}

// generates a JS code that runs the user's code against all test cases
func buildJsCode(ans models.Answer, q models.Question) (string, error) {

	codeToExec := ans.Code + `
	
	const inputs = [` + q.TestCases[0].Input

	for _, t := range q.TestCases[1:] {
		codeToExec += ` ,` + t.Input
	}

	codeToExec += `]
	const outputs = [` + q.TestCases[0].Output

	for _, t := range q.TestCases[1:] {
		codeToExec += ` ,` + t.Output
	}

	codeToExec += `]
	for (let i = 0; i < inputs.length; i++) {
        const ans = funcUser(inputs[i]);
        const json1 = JSON.stringify(ans);
        const json2 = JSON.stringify(outputs[i]);

        if (json1 !== json2) {
            return false;
        }
    }

    return true;`

	return codeToExec, nil
}

package services

import (
	"fmt"
	"leetcoder/models"
	"strings"
)

// build code for testing answer with all test cases
func buildTestCode(ans models.Answer, q models.Question) (string, error) {

	if ans.Lang == "python" {
		return buildPythonTest(ans, q)
	}

	if ans.Lang == "javascript" || ans.Lang == "js" {
		return buildJsTest(ans, q)
	}
	return "", fmt.Errorf("code language is not supported")
}

// generates a python code that runs the user's code against all test cases
func buildPythonTest(ans models.Answer, q models.Question) (string, error) {

	codeToExec := "import json\\n" + ans.Code + "\\ndef main():\\n\\tinputs = [" + q.TestCases[0].Input

	for _, t := range q.TestCases[1:] {
		codeToExec += " ," + t.Input
	}

	codeToExec += "]\\n\\toutputs = [" + q.TestCases[0].Output
	for _, t := range q.TestCases[1:] {
		codeToExec += " ," + t.Output
	}

	codeToExec += "]\\n\\tfor i in range(len(inputs)):\\n\\t\\tans = " + q.Name +
		"(inputs[i])\\n\\t\\tjson1 = json.dumps(ans)\\n\\t\\tjson2 = json.dumps(outputs[i])\\n\\t\\tif json1 != json2:\\n\\t\\t\\treturn False \\n\\treturn True \\n\\n\\nif __name__ == '__main__':\\n\\tprint(main())"

	return codeToExec, nil
}

// generates a JS code that runs the user's code against all test cases
func buildJsTest(ans models.Answer, q models.Question) (string, error) {

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

// build dockerfile for testing the answer
func buildDockerfile(code, lang string) (string, error) {

	if lang == "python" {
		return buildPythonDocker(code), nil
	}

	// if ans.Lang == "javascript" || ans.Lang == "js" {
	// 	return buildJsCode(ans, q)
	// }
	return "", fmt.Errorf("code language is not supported")
}

// build dockerfile for testing the answer - for python
func buildPythonDocker(pythonCode string) string {

	dockerfileContent := fmt.Sprintf(`
FROM python:3
WORKDIR /app
RUN echo '%s' > script.py
CMD ["python", "script.py"]
`, strings.ReplaceAll(pythonCode, "'", `'"'"'`))

	return dockerfileContent
}

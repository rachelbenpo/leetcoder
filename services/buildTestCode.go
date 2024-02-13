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
	return "", fmt.Errorf("code language is not supported: ", ans.Lang)
}

// generate a python code that runs the user's code against all test cases
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

// generate a JS code that runs the user's code against all test cases
func buildJsTest(ans models.Answer, q models.Question) (string, error) {

	codeToExec := ans.Code + "\\nfunction testAns(){\\n \\t let inputs = [" + q.TestCases[0].Input

	for _, t := range q.TestCases[1:] {
		codeToExec += " ," + t.Input
	}

	codeToExec += "]\\n \\t let outputs = [" + q.TestCases[0].Output

	for _, t := range q.TestCases[1:] {
		codeToExec += " ," + t.Output
	}

	codeToExec += "]\\n\\t for (let i = 0; i < inputs.length; i++) {\\n\\t\\t let ans =" + q.Name +
		"(inputs[i]);\\n\\t\\t let json1 = JSON.stringify(ans);\\n\\t\\t let json2 = JSON.stringify(outputs[i]);\\n\\t\\t if (json1 !== json2)\\n\\t\\t\\t return false; \\n\\t}\\n\\t return true; \\n} \\nconsole.log(testAns());"

	return codeToExec, nil
}

// build dockerfile for testing the answer
func buildDockerfile(code, lang string) (string, error) {

	if lang == "python" {
		return buildPythonDocker(code), nil
	}

	if lang == "javascript" || lang == "js" {
		return buildJSDocker(code), nil
	}
	return "", fmt.Errorf("code language is not supported")
}

// build dockerfile for running python code
func buildPythonDocker(pythonCode string) string {

	dockerfileContent := fmt.Sprintf(`
FROM python:3
WORKDIR /app
RUN echo '%s' > script.py
CMD ["python", "script.py"]
`, strings.ReplaceAll(pythonCode, "'", `'"'"'`))

	return dockerfileContent
}

// build a Dockerfile for running JavaScript code
func buildJSDocker(jsCode string) string {
	dockerfileContent := fmt.Sprintf(`
FROM node:14
WORKDIR /app
RUN echo '%s' > script.js
CMD ["node", "script.js"]
`, strings.ReplaceAll(jsCode, "'", `'"'"'`))

	return dockerfileContent
}

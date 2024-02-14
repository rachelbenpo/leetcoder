package clientCLI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "http://localhost:8080"

func GetAllQuestions() ([]byte, error) {
	url := baseURL + "/questions"
	resp, err := http.Get(url)
	return handleResponse(resp, err)
}

func GetQuestionByID(id int) ([]byte, error) {
	url := baseURL + "/questions/" + string(id)
	resp, err := http.Get(url)
	return handleResponse(resp, err)
}

func CreateQuestion(questionData map[string]interface{}) ([]byte, error) {

	url := baseURL + "/questions"
	jsonData, err := json.Marshal(questionData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	return handleResponse(resp, err)
}

func UpdateQuestion(id string, updatedQuestion map[string]interface{}) ([]byte, error) {

	url := baseURL + "/questions/" + id
	jsonData, err := json.Marshal(updatedQuestion)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	return handleResponse(resp, err)
}

func DeleteQuestion(id int) ([]byte, error) {

	url := baseURL + "/questions/" + string(id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	return handleResponse(resp, err)
}

func CheckAnswer(questionID string, answerData map[string]interface{}) {
	
	url := baseURL + "/questions/check-answer/" + questionID

	jsonData, err := json.Marshal(answerData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	handleResponse(resp, err)
}

// handle the http response
func handleResponse(resp *http.Response, err error) ([]byte, error) {

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	fmt.Println("Response:", resp.Status)
	fmt.Println(string(body))

	return body, nil
}

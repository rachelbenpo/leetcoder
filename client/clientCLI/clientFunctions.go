package clientCLI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"leetcode-client/models"
)

const baseURL = "http://localhost:8080"

func GetAllQuestions() ([]models.Question, error) {

	// send HTTP request to server
	url := baseURL + "/questions"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// convert the response body to questions array
	var questions []models.Question
	err = json.NewDecoder(resp.Body).Decode(&questions)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return nil, err
	}

	return questions, nil
}

func GetQuestionById(id string) (models.Question, error) {

	// send HTTP request
	url := baseURL + "/questions/" + id
	resp, err := http.Get(url)
	if err != nil {
		return models.Question{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Question{}, fmt.Errorf("Server returned non-OK status: %s", resp.Status)
	}

	// convert response body to question
	var question models.Question
	err = json.NewDecoder(resp.Body).Decode(&question)
	if err != nil {
		return models.Question{}, fmt.Errorf("Error decoding response body: %s", err)
	}

	return question, nil
}

func CreateQuestion(question models.Question) (int, error) {

	url := baseURL + "/questions"

	// convert input question to json
	jsonData, err := json.Marshal(question)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return 0, err
	}

	// send HTTP request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return 0, err
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Error: Unexpected status code", resp.StatusCode)
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// extract data from response
	var id int
	err = json.NewDecoder(resp.Body).Decode(&id)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return 0, err
	}

	return id, nil
}

func UpdateQuestion(id string, question models.Question) error {

	// send HTTP request
	url := baseURL + "/questions/" + id
	jsonData, err := json.Marshal(question)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making PUT request:", err)
		return err
	}
	defer resp.Body.Close()

	// check response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code", resp.StatusCode)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func DeleteQuestion(id string) error {

	url := baseURL + "/questions/" + id

	// send HTTP request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making DELETE request:", err)
		return err
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code", resp.StatusCode)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func CheckAnswer(id string, ans models.Answer) (string, error) {

	url := baseURL + "/questions/check-answer/" + id

	// convert input answer to json
	jsonData, err := json.Marshal(ans)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return "", err
	}

	// send HTTP request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code", resp.StatusCode)
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// get data from response body
	var correct string
	err = json.NewDecoder(resp.Body).Decode(&correct)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return "", err
	}

	return correct, nil
}

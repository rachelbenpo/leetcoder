package clientCLI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "http://localhost:8080"

func GetAllQuestions() ([]Question, error) {

	// send HTTP request to server
	url := baseURL + "/questions"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ERROR: status code: %d message: %d", resp.StatusCode, resp.Body)
	}

	// convert the response body to questions array
	var questions []Question
	err = json.NewDecoder(resp.Body).Decode(&questions)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return nil, err
	}

	return questions, nil
}

func GetQuestionById(id string) (Question, error) {

	// send HTTP request
	url := baseURL + "/questions/" + id
	resp, err := http.Get(url)
	if err != nil {
		return Question{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Question{}, fmt.Errorf("server returned error. \nstatus code: %d \nmessage: %d", resp.StatusCode, resp.Body)
	}

	// convert response body to question
	var question Question
	err = json.NewDecoder(resp.Body).Decode(&question)
	if err != nil {
		return Question{}, fmt.Errorf("error decoding response body: %s", err)
	}

	return question, nil
}

func CreateQuestion(question Question) (int, error) {

	url := baseURL + "/questions"

	// convert input question to json
	jsonData, err := json.Marshal(question)
	if err != nil {
		return 0, fmt.Errorf("error marshalling JSON:", err)
	}

	// send HTTP request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("server returned error. \nstatus code: %d \nmessage: %d", resp.StatusCode, resp.Body)
	}

	// extract data from response
	var r struct {
		id      int
		message string
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return 0, fmt.Errorf("Error decoding response: ", err)
	}

	return r.id, nil
}

func UpdateQuestion(id string, question Question) error {

	// send HTTP request
	url := baseURL + "/questions/" + id
	jsonData, err := json.Marshal(question)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: ", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error creating request: ", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error making PUT request: ", err)
	}
	defer resp.Body.Close()

	// check response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned error. \nstatus code: %d \nmessage: %d", resp.StatusCode, resp.Body)
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

func CheckAnswer(id string, ans Answer) (string, error) {

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
		return "", fmt.Errorf("error: %d", resp.StatusCode)
	}

	// get data from response body

	// var isCorrect struct {
	// 	correct string
	// }

	is := map[string]string{}

	err = json.NewDecoder(resp.Body).Decode(&is)
	if err != nil {
		fmt.Println("Errorrrr decoding response:", err)
		return "", err
	}

	return is["correct"], nil
}

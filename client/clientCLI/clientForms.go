package clientCLI

import (
	// "encoding/json"
	"fmt"
	"os"
	// "strconv"

	"github.com/charmbracelet/huh"
	// "github.com/charmbracelet/lipgloss"
	"leetcode-client/models"
)

func createForm() {

	// help variables for storing form data
	var question = models.Question{
		ID:           0,
		Name:         "",
		Instructions: "",
		Answer:       "",
		TestCases:    []models.TestCase{},
	}

	var ans = models.Answer{
		Lang: "",
		Code: "",
	}

	var id string
	var testCasesCount string
	var operation string

	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions("Get All Questions", "Get Question by ID", "Create Question", "Update Question", "Delete Question", "Check Answer")...).
				Title("Choose Operation").
				Value(&operation),
		))

	err := form.Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	createQuestionForm := huh.NewForm(
		// create a new question form
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Placeholder("Enter Name").
				Value(&question.Name),
			huh.NewInput().
				Title("Instructions").
				Placeholder("Enter Instructions").
				Value(&question.Instructions),
			huh.NewInput().
				Title("Answer - optional").
				Placeholder("Enter Answer if you want").
				Value(&question.Answer),
			huh.NewInput().
				Title("num of test cases").
				Value(&testCasesCount),

			huh.NewMultiSelect[models.TestCase]().
				Title("Test Cases").
				Options(huh.NewOptions(models.TestCase{ID: 1, Input: "Input1", Output: "Output1"}, models.TestCase{ID: 2, Input: "Input2", Output: "Output2"})...).
				Filterable(true),
		))

	getIDForm := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Question ID").
			Placeholder("Enter Question ID").
			Value(&id),
	))

	checkAnswerForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Language").
				Placeholder("Enter Language").
				Value(&ans.Lang),
			huh.NewInput().
				Title("Code").
				Placeholder("Enter Code").
				Value(new(string)),
		))

	for true {

		switch operation {

		case "Get All Questions":
			questions, err := GetAllQuestions()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			printQuestions(questions)

		case "Get Question by ID":

			err := getIDForm.Run()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			q, err := GetQuestionById(id)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			printQuestion(q)

		case "Create Question":

			err := createQuestionForm.Run()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			id, err := CreateQuestion(question)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Question created with ID:", id)

		case "Update Question":

			err := getIDForm.Run()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			err = createQuestionForm.Run()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			err = UpdateQuestion(id, question)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Question updated successfully.")

		case "Delete Question":
			err := getIDForm.Run()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			err = DeleteQuestion(id)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Question deleted successfully.")

		case "Check Answer":
			err := checkAnswerForm.Run()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			correct, err := CheckAnswer(id, ans)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Answer is correct:", correct)

		default:
			fmt.Println("Invalid operation selected.")
		}
	}
}

func printQuestions(questions []models.Question) {
	for _, q := range questions {
		printQuestion(q)
	}
}

func printQuestion(question models.Question) {
	fmt.Printf("ID: %d\nName: %s\nInstructions: %s", question.ID, question.Name, question.Instructions)

	fmt.Println("\n-------------------------------")
}

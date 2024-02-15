package clientCLI

import (
	"fmt"
	"strconv"

	"leetcode-client/models"

	"github.com/charmbracelet/huh"
)

// help variables for storing form data
var (
	question      = models.Question{}
	ans           = models.Answer{}
	testCase      = models.TestCase{}
	shortQuestion = models.QuestionShort{}

	id             string
	testCasesCount string
	operation      string = "Choose operation"
	wantTOAnswer   bool   = false

	questionsList  []models.Question
	shortQuestions []models.QuestionShort
)

// run the CLI
func RunCLIforms() {

	var err error

	for {
		switch operation {

		case "Choose operation":
			if err = chooseOperation(); err != nil {
				fmt.Println("Error:", err)
				return
			}

		case "Get All Questions":
			if err = runGetAllForm(); err != nil {
				fmt.Println("error getting all questions from server: ", err)
				operation = "Choose operation"
			}

		case "Get Question by ID":
			if err = RunGetByIdForm(); err != nil {
				fmt.Println("error getting question from server: ", err)
				operation = "Choose operation"
			}

		case "Show single question":
			if err = runShowSingleForm(); err != nil {
				fmt.Println("error getting question from server: ", err)
				operation = "Choose operation"
			}

		case "Create Question":
			if err = runCreateQuestionForm(); err != nil {
				fmt.Println("error creating question: ", err)
			}
			operation = "Choose operation"

		case "Update Question":
			if err := runUpdateForm(); err != nil {
				fmt.Println("error updating question: ", err)
			}
			operation = "Choose operation"

		case "Delete Question":
			if err = runDeleteForm(); err != nil {
				fmt.Println("error deleting question:", err)
			}
			operation = "Choose operation"

		case "Check Answer":
			if err = runCheckForm(); err != nil {
				operation = "Choose operation"
				fmt.Println("error testing your answer: ", err)
			}
		case "Exit":
			return

		default:
			fmt.Println("Invalid operation selected.")
			return
		}
	}
}

// run form for showing all questions and choosing once
func runGetAllForm() error {

	// get all questions from the server
	questionsList, err := GetAllQuestions()
	if err != nil {
		return err
	}

	shortQuestions = getShortQuestions(questionsList)

	// show all questions and choose once
	allQuestionsForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[models.QuestionShort]().
				Options(huh.NewOptions(shortQuestions...)...).
				Title("Choose Question").
				Value(&shortQuestion)),
	)

	if err = allQuestionsForm.Run(); err != nil {
		return err
	}
	operation = "Show single question"
	id = fmt.Sprint(shortQuestion.ID)
	return nil
}

// run form for get question by id
func RunGetByIdForm() error {

	// form for enter id
	var getIDForm = huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Question ID").
			Placeholder("Enter Question ID").
			Value(&id),
	))
	getIDForm.View()
	err := getIDForm.Run()
	if err != nil {
		return err
	}

	if question, err = GetQuestionById(id); err != nil {
		return err
	}

	// form for showing a single question
	singleQuestionForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(question.Name).
				Description(question.Instructions),

			huh.NewConfirm().
				Affirmative("try to answer").
				Negative("cancel").
				Value(&wantTOAnswer),
		),
	)

	if err = singleQuestionForm.Run(); err != nil {
		return err
	}
	if wantTOAnswer {
		operation = "Check Answer"
	} else {
		operation = "Choose operation"
	}
	return nil
}

// run form for showing a single question
func runShowSingleForm() error {

	var err error

	if question, err = GetQuestionById(id); err != nil {
		return err
	}

	// form for show a single question
	singleQuestionForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(question.Name).
				Description(question.Instructions),

			huh.NewConfirm().
				Affirmative("try to answer").
				Negative("cancel").
				Value(&wantTOAnswer),
		),
	)

	if err = singleQuestionForm.Run(); err != nil {
		return err
	}
	if wantTOAnswer {
		operation = "Check Answer"
	} else {
		operation = "Choose operation"
	}
	return nil
}

// run form for deletion
func runDeleteForm() error {

	// form for enter id
	getIDForm := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Question ID").
			Placeholder("Enter Question ID").
			Value(&id),
	))

	if err := getIDForm.Run(); err != nil {
		return err
	}

	if err := DeleteQuestion(id); err != nil {
		return err
	}
	return nil
}

// run form for update a question
func runUpdateForm() error {

	// form for enter id
	var getIDForm = huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Question ID").
			Placeholder("Enter Question ID").
			Value(&id),
	))

	if err := getIDForm.Run(); err != nil {
		return err
	}

	updateQuestionForm := huh.NewForm(
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
		).Title("fill only the parameters you want to update"),
	)

	err := updateQuestionForm.Run()
	if err != nil {
		return err
	}
	if testCasesCount == "" {
		question.TestCases = nil
	} else {
		count, err := strconv.Atoi(testCasesCount)
		if err != nil {
			return err
		}

		//run form for creating test cases
		for i := 0; i < count; i++ {
			testCaseForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Input").
						Placeholder("Enter Input").
						Value(&testCase.Input),
					huh.NewInput().
						Title("Output").
						Placeholder("Enter Output").
						Value(&testCase.Output),
				),
			)
			if err := testCaseForm.Run(); err != nil {
				return err
			}
			question.TestCases = append(question.TestCases, testCase)
			testCase = models.TestCase{Input: "", Output: ""}
		}
	}

	// send HTTP request to the server
	if err := UpdateQuestion(id, question); err != nil {
		return err
	}
	return nil
}

// run form for creating a new question
func runCreateQuestionForm() error {

	createQuestionForm := huh.NewForm(
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
		),
	)

	err := createQuestionForm.Run()
	if err != nil {
		return err
	}

	count, err := strconv.Atoi(testCasesCount)
	if err != nil {
		return err
	}

	//run form for creating test cases
	for i := 0; i < count; i++ {
		testCaseForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Input").
					Placeholder("Enter Input").
					Value(&testCase.Input),
				huh.NewInput().
					Title("Output").
					Placeholder("Enter Output").
					Value(&testCase.Output),
			),
		)
		err := testCaseForm.Run()
		if err != nil {
			return err
		}
		question.TestCases = append(question.TestCases, testCase)
		testCase = models.TestCase{Input: "", Output: ""}
	}

	// send HTTP request to the server
	_, err = CreateQuestion(question)
	if err != nil {
		return err
	}
	return nil
}

// run form for checking an answer
func runCheckForm() error {

	if id == "" || id == "0" {
		// run form for enter id
		getIDForm := huh.NewForm(huh.NewGroup(
			huh.NewInput().
				Title("Question ID").
				Placeholder("Enter Question ID").
				Value(&id),
		))

		if err := getIDForm.Run(); err != nil {
			return err
		}
	}

	checkAnswerForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOption("javascript", "javascript"), huh.NewOption("python", "python")).
				Title("Language").
				Value(&ans.Lang),
			huh.NewText().
				Title("Code").
				Placeholder("Enter Code").
				Value(&ans.Code),
		),
	)

	if err := checkAnswerForm.Run(); err != nil {
		return err
	}

	// send question to server for testing
	correct, err := CheckAnswer(id, ans)
	if err != nil {
		return err
	}

	// check response
	if correct == "true" {
		fmt.Println("correct answer.")
		operation = "Choose operation"
	} else {
		fmt.Println("wrong answer.")
		fmt.Println("message: ", correct)
		operation = "Check Answer"
		ans = models.Answer{Lang: "", Code: ""}
	}
	return nil
}

// run form for choosing which operation to perform
func chooseOperation() error {

	initVariables()

	operationForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions("Get All Questions", "Get Question by ID", "Create Question", "Update Question", "Delete Question", "Check Answer", "Exit")...).
				Title("Choose Operation").
				Value(&operation),
		),
	)
	if err := operationForm.Run(); err != nil {
		return err
	}
	return nil
}

// init help variables
func initVariables() {

	question = models.Question{}
	ans = models.Answer{}
	testCase = models.TestCase{}
	shortQuestion = models.QuestionShort{}

	id = ""
	testCasesCount = ""

	wantTOAnswer = false

	questionsList = []models.Question{}
	shortQuestions = []models.QuestionShort{}
}

// get array of short questions
func getShortQuestions(questions []models.Question) []models.QuestionShort {

	var shortQuestions []models.QuestionShort
	for _, q := range questions {
		shortQuestions = append(shortQuestions, models.QuestionShort{ID: q.ID, Name: q.Name})
	}
	return shortQuestions
}

package clientCLI

import (
	"fmt"
	"strconv"

	"leetcode-client/models"

	"github.com/charmbracelet/huh"
)

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

var testCase = models.TestCase{
	Input:  "",
	Output: "",
}

var (
	id             string
	testCasesCount string
	operation      string = "choose operation"
	shortQuestion  models.QuestionShort
	wantTOAnswer   bool = false

	questionsList  []models.Question
	shortQuestions []models.QuestionShort
	names          []string
	IDs            []int
)

/////////////////////////////////

// forms:

// form for choosing which operation to perform
var operationForm = huh.NewForm(
	huh.NewGroup(
		huh.NewSelect[string]().
			Options(huh.NewOptions("Get All Questions", "Get Question by ID", "Create Question", "Update Question", "Delete Question", "Check Answer")...).
			Title("Choose Operation").
			Value(&operation),
	),
)

// form for creating a new question (or updating an existing question)
var createQuestionForm = huh.NewForm(
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

// form for creating test cases
var testCaseForm = huh.NewForm(
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

// form for enter id
var getIDForm = huh.NewForm(huh.NewGroup(
	huh.NewInput().
		Title("Question ID").
		Placeholder("Enter Question ID").
		Value(&id),
),
)

// form for checking an answer
var checkAnswerForm = huh.NewForm(
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

// form for showing all questions and choosing once
var allQuestionsForm = huh.NewForm(
	huh.NewGroup(
		huh.NewSelect[models.QuestionShort]().
			Options(huh.NewOptions(shortQuestions...)...).
			Title("Choose Question").
			Value(&shortQuestion)),
)

// form for showing a single question
var singleQuestionForm = huh.NewForm(
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

///////////////////////////////////////

func CreateForm() {

	var err error
	for {
		// choose which operation to perform
		switch operation {
		case "choose operation":
			if err = operationForm.Run(); err != nil {
				fmt.Println("Error:", err)
				return
			}

		case "Get All Questions":
			if err = getAll(); err != nil {
				fmt.Println("error getting all questions from server: ", err)
			}

		case "Get Question by ID":
			if err = getById(); err != nil {
				fmt.Println("error getting question from server: ", err)
			}
		case "show single question":
			if err = showSingle(); err != nil {
				fmt.Println("error getting question from server: ", err)
			}

		case "Create Question":
			if err = create(); err != nil {
				fmt.Println("error creating question: ", err)
			}

		case "Update Question":
			if err := update(); err != nil {
				fmt.Println("error updating question: ", err)
			}

		case "Delete Question":
			if err = delete(); err != nil {
				fmt.Println("error deleting question:", err)
			}

		case "Check Answer":
			if err = check(); err != nil {
				fmt.Println("error testing your answer: ", err)
			}
		default:
			fmt.Println("Invalid operation selected.")
			return
		}
	}
}

// ////////////////////////
func printQuestions(questions []models.Question) {
	for _, q := range questions {
		printQuestion(q)
	}
}

func printQuestion(question models.Question) {
	fmt.Printf("ID: %d\nName: %s\nInstructions: %s", question.ID, question.Name, question.Instructions)

	fmt.Println("\n-------------------------------")
}

// get array of all questions Ids
func getIds(questions []models.Question) []string {
	var ids []string
	for _, q := range questions {
		ids = append(ids, string(q.ID))
	}
	return ids
}

// get array of short questions
func getShortQuestions(questions []models.Question) []models.QuestionShort {

	var shortQuestions []models.QuestionShort
	for _, q := range questions {
		shortQuestions = append(shortQuestions, models.QuestionShort{ID: q.ID, Name: q.Name})
	}
	return shortQuestions
}

func getNames(questions []models.Question) []string {

	var names []string
	for _, q := range questions {
		names = append(names, q.Name)
	}
	return names
}

///////////////////////////

func getAll() error {

	// get all questions from the server
	questionsList, err := GetAllQuestions()
	if err != nil {
		return err
	}

	shortQuestions = getShortQuestions(questionsList)

	// show all questions and choose once
	allQuestionsForm = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[models.QuestionShort]().
				Options(huh.NewOptions(shortQuestions...)...).
				Title("Choose Question").
				Value(&shortQuestion)),
	)

	if err = allQuestionsForm.Run(); err != nil {
		return err
	}
	operation = "Get Question by ID"
	id = fmt.Sprint(shortQuestion.ID)
	return nil
}

func delete() error {

	if err := getIDForm.Run(); err != nil {
		return err
	}

	if err := DeleteQuestion(id); err != nil {
		return err
	}
	return nil
}

func update() error {

	if err := getIDForm.Run(); err != nil {
		return err
	}

	if err := createQuestionForm.Run(); err != nil {
		return err
	}

	if err := UpdateQuestion(id, question); err != nil {
		return err
	}
	return nil
}

func getById() error {

	err := getIDForm.Run()
	if err != nil {
		return err
	}

	if question, err = GetQuestionById(id); err != nil {
		return err
	}

	// form for showing a single question
	singleQuestionForm = huh.NewForm(
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
		operation = "choose operation"
	}
	return nil
}

func showSingle() error {
	var err error

	if question, err = GetQuestionById(id); err != nil {
		return err
	}

	// form for showing a single question
	singleQuestionForm = huh.NewForm(
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
		operation = "choose operation"
	}
	return nil
}

func create() error {

	// run form of creating a new question
	err := createQuestionForm.Run()
	if err != nil {
		return err
	}

	count, err := strconv.Atoi(testCasesCount)
	if err != nil {
		return err
	}

	// create test cases
	for i := 0; i < count; i++ {
		err := testCaseForm.Run()
		if err != nil {
			return err
		}
		question.TestCases = append(question.TestCases, testCase)
	}

	_, err = CreateQuestion(question)
	if err != nil {
		return err
	}
	return nil
}

func check() error {

	// run form of checking an answer
	if err := checkAnswerForm.Run(); err != nil {
		return err
	}

	// send question to server for testing
	correct, err := CheckAnswer(id, ans)
	if err != nil {
		return err
	}

	// print response
	if correct == "true" {
		fmt.Println("correct answer.")
	} else {
		fmt.Println("wrong answer.")
		fmt.Println("message: ", correct)
	}
	return nil
}

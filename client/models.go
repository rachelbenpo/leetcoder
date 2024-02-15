package main

type Question struct {
	ID           int
	Name         string
	Instructions string
	Answer       string
	TestCases    []TestCase `json:"test_cases"`
}

type TestCase struct {
	ID     int
	Input  string
	Output string
}

type Answer struct {
	Lang string
	Code string
}

type QuestionShort struct {
	ID   int
	Name string
}

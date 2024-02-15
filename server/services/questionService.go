package services

import (
	"database/sql"
	"fmt"
	"leetcode-server/config"
	"leetcode-server/models"
)

// services - logic

func GetAllQuestions() ([]models.Question, error) {

	// connect to DB
	db, err := sql.Open("mysql", config.DBConnectionString)
	if err != nil {
		return []models.Question{}, err
	}
	defer db.Close()

	// exec SQL
	rows, err := db.Query("SELECT q.*, t.id, t.input, t.output FROM questions q LEFT JOIN test_cases t ON q.id = t.question_id")
	if err != nil {
		return []models.Question{}, err
	}
	defer rows.Close()

	// get the data
	var questionsMap = make(map[int]models.Question)

	for rows.Next() {
		var question models.Question
		question.TestCases = []models.TestCase{}
		var testCaseID sql.NullInt64
		var input, output sql.NullString

		err := rows.Scan(&question.ID, &question.Name, &question.Instructions, &question.Answer, &testCaseID, &input, &output)
		if err != nil {
			return []models.Question{}, err
		}

		q, ok := questionsMap[question.ID]

		// add the new Question if doesn't exist
		if !ok {
			q = question
		}

		// Append test case
		if testCaseID.Valid {
			q.TestCases = append(q.TestCases, models.TestCase{
				ID:     int(testCaseID.Int64),
				Input:  input.String,
				Output: output.String,
			})
		}

		questionsMap[question.ID] = q
	}

	// Convert map to question array
	var questions []models.Question
	for _, question := range questionsMap {
		questions = append(questions, question)
	}

	return questions, nil
}

func GetQuestionById(id string) (models.Question, error) {

	// connect to DB
	db, err := sql.Open("mysql", config.DBConnectionString)
	if err != nil {
		return models.Question{}, err
	}
	defer db.Close()

	// exec SQL
	rows, err := db.Query("SELECT q.*, t.id, t.input, t.output FROM questions q LEFT JOIN test_cases t ON q.id = t.question_id WHERE q.id = ?", id)
	if err != nil {
		return models.Question{}, err
	}
	defer rows.Close()

	// get the data
	var question models.Question
	question.TestCases = []models.TestCase{}

	for rows.Next() {
		var testCaseID sql.NullInt64
		var input, output sql.NullString

		err := rows.Scan(&question.ID, &question.Name, &question.Instructions, &question.Answer, &testCaseID, &input, &output)
		if err != nil {
			return models.Question{}, err
		}

		// Append test case if valid
		if testCaseID.Valid {
			question.TestCases = append(question.TestCases, models.TestCase{
				ID:     int(testCaseID.Int64),
				Input:  input.String,
				Output: output.String,
			})
		}
	}

	// if no question is found
	if question.ID == 0 {
		return models.Question{}, fmt.Errorf("question not found")
	}

	return question, nil
}

func CreateQuestion(question models.Question) (int, error) {

	// connect to DB
	db, err := sql.Open("mysql", config.DBConnectionString)
	if err != nil {
		return 0, err
	}

	// exec SQL
	stmt, err := db.Prepare("INSERT INTO questions (name, instructions, answer) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(question.Name, question.Instructions, question.Answer)
	if err != nil {
		return 0, err
	}

	// get the new id
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// save testcases
	err = createTestCases(db, question.TestCases, int(id))
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func UpdateQuestion(id string, updatedQuestion models.Question) error {

	// connect to DB
	db, err := sql.Open("mysql", config.DBConnectionString)
	if err != nil {
		return err
	}

	// check if question exists
	row := db.QueryRow("SELECT * FROM questions WHERE id = ?", id)
	var question models.Question
	err = row.Scan(&question.ID, &question.Name, &question.Instructions, &question.Answer)
	if err != nil {
		return err
	}

	// update values
	if updatedQuestion.Name != "" {
		question.Name = updatedQuestion.Name
	}
	if updatedQuestion.Instructions != "" {
		question.Instructions = updatedQuestion.Instructions
	}
	if updatedQuestion.Answer != "" {
		question.Answer = updatedQuestion.Answer
	}

	// update test cases
	if updatedQuestion.TestCases != nil {
		err := updateTestCases(db, updatedQuestion.TestCases, question.ID)
		if err != nil {
			return err
		}
	}

	// exec SQL
	stmt, err := db.Prepare("UPDATE questions SET name = ?, instructions = ?, answer = ? WHERE id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(question.Name, question.Instructions, question.Answer, id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteQuestion(id string) error {

	// connect to DB
	db, err := sql.Open("mysql", config.DBConnectionString)
	if err != nil {
		return err
	}

	// check if question exists
	row := db.QueryRow("SELECT * FROM questions WHERE id = ?", id)
	var question models.Question
	err = row.Scan(&question.ID, &question.Name, &question.Instructions, &question.Answer)
	if err != nil {
		return err
	}

	// Delete test cases
	_, err = db.Exec("DELETE FROM test_cases WHERE question_id = ?", question.ID)
	if err != nil {
		return err
	}

	// delete question
	_, err = db.Exec("DELETE FROM questions WHERE id = ?", question.ID)
	if err != nil {
		return err
	}
	return nil
}

func createTestCases(db *sql.DB, testCases []models.TestCase, questionId int) error {

	stmt, err := db.Prepare("INSERT INTO test_cases (question_id, input, output) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, testCase := range testCases {
		_, err := stmt.Exec(questionId, testCase.Input, testCase.Output)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateTestCases(db *sql.DB, testCases []models.TestCase, questionID int) error {

	// Delete existing test cases
	_, err := db.Exec("DELETE FROM test_cases WHERE question_id = ?", questionID)
	if err != nil {
		return err
	}

	// save new testcases
	return createTestCases(db, testCases, questionID)
}

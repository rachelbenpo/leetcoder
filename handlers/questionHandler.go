package handlers

import (
	"fmt"
	"leetcoder/services"
	"leetcoder/models"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// CRUD functions

func GetAllQuestions(c *gin.Context) {

	questions, err := services.GetAllQuestions()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	c.JSON(200, questions)
}

func GetQuestionById(c *gin.Context) {

	id := c.Param("id")

	question, err := services.GetQuestionById(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	c.JSON(200, question)
}

func CreateQuestion(c *gin.Context) {

	// get question params
	var question models.Question
	err := c.BindJSON(&question)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	id, err := services.CreateQuestion(question)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	c.JSON(201, gin.H{"id": id, "message": "question created"})
}

func UpdateQuestion(c *gin.Context) {

	id := c.Param("id")

	// get question params
	var updatedQuestion models.Question

	err := c.BindJSON(&updatedQuestion)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	err = services.UpdateQuestion(id, updatedQuestion)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	c.JSON(200, gin.H{"message": "question updated"})
}

func DeleteQuestion(c *gin.Context) {

	id := c.Param("id")

	// delete question
	err := services.DeleteQuestion(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	c.JSON(200, gin.H{"message": "question deleted"})
}

/*
func getAllQuestions(c *gin.Context) {

	// connect to DB
	db, err := sql.Open("mysql", DBConnectionString)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}
	defer db.Close()

	// exec SQL
	rows, err := db.Query("SELECT q.*, t.id, t.input, t.output FROM questions q LEFT JOIN test_cases t ON q.id = t.question_id")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}
	defer rows.Close()

	// get the data
	var questionsMap = make(map[int]Question)

	for rows.Next() {
		var question Question
		question.TestCases = []TestCase{}
		var testCaseID sql.NullInt64
		var input, output sql.NullString

		err := rows.Scan(&question.ID, &question.Name, &question.Instructions, &question.Answer, &testCaseID, &input, &output)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			fmt.Print(err)
			return
		}

		q, ok := questionsMap[question.ID]

		// add the new Question if doesn't exist
		if !ok {
			q = question
		}

		// Append test case
		if testCaseID.Valid {
			q.TestCases = append(q.TestCases, TestCase{
				ID:     int(testCaseID.Int64),
				Input:  input.String,
				Output: output.String,
			})
		}

		questionsMap[question.ID] = q
	}

	// Convert map to question array
	var questions []Question
	for _, question := range questionsMap {
		questions = append(questions, question)
	}

	c.JSON(200, questions)
}

func getQuestionById(c *gin.Context) {

	id := c.Param("id")

	// connect to DB
	db, err := sql.Open("mysql", DBConnectionString)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}
	defer db.Close()

	// exec SQL
	rows, err := db.Query("SELECT q.*, t.id, t.input, t.output FROM questions q LEFT JOIN test_cases t ON q.id = t.question_id WHERE q.id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}
	defer rows.Close()

	// get the data
	var question Question
	question.TestCases = []TestCase{}

	for rows.Next() {
		var testCaseID sql.NullInt64
		var input, output sql.NullString

		err := rows.Scan(&question.ID, &question.Name, &question.Instructions, &question.Answer, &testCaseID, &input, &output)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			fmt.Print(err)
			return
		}

		// Append test case if valid
		if testCaseID.Valid {
			question.TestCases = append(question.TestCases, TestCase{
				ID:     int(testCaseID.Int64),
				Input:  input.String,
				Output: output.String,
			})
		}
	}

	// if no question is found
	if question.ID == 0 {
		c.JSON(404, gin.H{"error": "question not found"})
		fmt.Print("error: question not found")
		return
	}

	c.JSON(200, question)
}

func createQuestion(c *gin.Context) {

	// get question params
	var question Question
	err := c.BindJSON(&question)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	// connect to DB
	db, err := sql.Open("mysql", DBConnectionString)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// exec SQL
	stmt, err := db.Prepare("INSERT INTO questions (name, instructions, answer) VALUES (?, ?, ?)")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(question.Name, question.Instructions, question.Answer)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	// get the new id
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	// save testcases
	err = createTestCases(db, question.TestCases, int(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	c.JSON(201, gin.H{"id": id, "message": "question created"})
}

func createTestCases(db *sql.DB, testCases []TestCase, questionId int) error {

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

func updateQuestion(c *gin.Context) {

	id := c.Param("id")

	// get question params
	var updatedQuestion Question

	err := c.BindJSON(&updatedQuestion)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	// connect to DB
	db, err := sql.Open("mysql", DBConnectionString)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	// check if question exists
	row := db.QueryRow("SELECT * FROM questions WHERE id = ?", id)
	var question Question
	err = row.Scan(&question.ID, &question.Name, &question.Instructions, &question.Answer)
	if err != nil {
		c.JSON(404, gin.H{"error": "question not found"})
		fmt.Print(err)
		return
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
			c.JSON(500, gin.H{"error": err.Error()})
			fmt.Print(err)
			return
		}
	}

	// exec SQL
	stmt, err := db.Prepare("UPDATE questions SET name = ?, instructions = ?, answer = ? WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(question.Name, question.Instructions, question.Answer, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	c.JSON(200, gin.H{"message": "question updated"})
}

func updateTestCases(db *sql.DB, testCases []TestCase, questionID int) error {

	// Delete existing test cases
	_, err := db.Exec("DELETE FROM test_cases WHERE question_id = ?", questionID)
	if err != nil {
		return err
	}

	// save new testcases
	return createTestCases(db, testCases, questionID)
}

func deleteQuestion(c *gin.Context) {

	id := c.Param("id")

	// connect to DB
	db, err := sql.Open("mysql", DBConnectionString)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	// check if question exists
	row := db.QueryRow("SELECT * FROM questions WHERE id = ?", id)
	var question Question
	err = row.Scan(&question.ID, &question.Name, &question.Instructions, &question.Answer)
	if err != nil {
		c.JSON(404, gin.H{"error": "question not found"})
		fmt.Print(err)
		return
	}

	// Delete test cases
	_, err = db.Exec("DELETE FROM test_cases WHERE question_id = ?", question.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	// delete question
	_, err = db.Exec("DELETE FROM questions WHERE id = ?", question.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	c.JSON(200, gin.H{"message": "question deleted"})
}
*/

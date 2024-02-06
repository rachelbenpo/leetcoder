package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type Question struct {
	ID           int
	Name         string
	Instructions string
	Answer       string
	TestCases    []*TestCase
}

type TestCase struct {
	ID     int
	Input  string
	Output string
}

func main() {

	fmt.Printf("server runs on localhost:8080\n")

	router := gin.Default()

	// set url for CRUD
	router.GET("/", getAllQuestions)
	router.GET("/questions", getAllQuestions)
	router.GET("/questions/:id", getQuestionById)
	router.POST("/questions", createQuestion)
	router.PUT("/questions/:id", updateQuestion)
	router.DELETE("/questions/:id", deleteQuestion)

	router.Run(":8080")
}

// CRUD functions

func getAllQuestions(c *gin.Context) {

	// connect to DB
	db, err := sql.Open("mysql", DBConnectionString)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	// exec SQL
	rows, err := db.Query("SELECT * FROM questions")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// get the data
	var questions []Question
	for rows.Next() {
		var question Question
		err := rows.Scan(&question.ID, &question.Name, &question.Instructions, &question.Answer)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		questions = append(questions, question)
	}

	c.JSON(200, questions)
}

func getQuestionById(c *gin.Context) {

	// get id and convert it to int
	idInput := c.Param("id")

	id, err := strconv.Atoi(idInput)
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

	// exec SQL
	row := db.QueryRow("SELECT * FROM questions WHERE id = ?", id)

	// scan and return data
	var question Question
	err = row.Scan(&question.ID, &question.Name, &question.Instructions, &question.Answer)

	if err != nil {
		c.JSON(404, gin.H{"error": "question not found"})
		fmt.Print(err)
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

func createTestCases(db *sql.DB, testCases []*TestCase, questionId int) error {

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

	// get question id
	id := c.Param("id")

	// get question params
	var updatedQuestion struct {
		ID			 *int		 `json:"id"`
		Name         *string     `json:"name"`
		Instructions *string     `json:"instructions"`
		Answer       *string     `json:"answer"`
		TestCases    []*TestCase `json:"test_cases"`
	}

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
	if updatedQuestion.Name != nil {
		question.Name = *updatedQuestion.Name
	}

	if updatedQuestion.Instructions != nil {
		question.Instructions = *updatedQuestion.Instructions
	}

	if updatedQuestion.Answer != nil {
		question.Answer = *updatedQuestion.Answer
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

/*func updateQuestion(c *gin.Context) {

	id := c.Param("id")


	var question Question
	err := c.BindJSON(&question)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	db, err := sql.Open("mysql", DBConnectionString)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	row := db.QueryRow("SELECT * FROM questions WHERE id = ?", id)
	err = row.Scan(&question.ID, &question.Name, &question.Instructions, &question.Answer)
	if err != nil {
		c.JSON(404, gin.H{"error": "question not found"})
		return
	}

	stmt, err := db.Prepare("UPDATE questions SET name = ?, instructions = ?, answer = ? WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(question.Name, question.Instructions, question.Answer, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = updateTestCases(db, question.TestCases, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "question updated"})
}*/

func updateTestCases(db *sql.DB, testCases []*TestCase, questionId int) error {

	// delete existing testcases
	stmt, err := db.Prepare("DELETE FROM test_cases WHERE question_id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(questionId)
	if err != nil {
		return err
	}

	// save new testcases
	return createTestCases(db, testCases, questionId)
}

func deleteQuestion(c *gin.Context) {

	id := c.Param("id")

	// connect to DB
	db, err := sql.Open("mysql", DBConnectionString)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// check if question exists
	row := db.QueryRow("SELECT * FROM questions WHERE id = ?", id)
	var question Question
	err = row.Scan(&question.ID, &question.Name, &question.Instructions, &question.Answer)
	if err != nil {
		c.JSON(404, gin.H{"error": "question not found"})
		return
	}

	// delete test cases
	stmt, err := db.Prepare("DELETE FROM test_cases WHERE question_id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(question.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// delete question
	stmt, err = db.Prepare("DELETE FROM questions WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "question deleted"})
}

// save connection string in safe place
// documentation, arrange, clear
// structure of project
// delete testcases function
// is all params come from user in update or not
// answer is need?

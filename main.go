package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// set url for CRUD
	router.GET("/", getAllQuestions)
	router.GET("/questions", getAllQuestions)
	router.GET("/questions/:id", getQuestionById)
	router.POST("/questions", createQuestion)
	router.PUT("/questions/:id", updateQuestion)
	router.DELETE("/questions/:id", deleteQuestion)
	router.POST("/questions/{id}/check-answer", checkAnswer)

	router.Run(":8080")

	fmt.Printf("server runs on localhost:8080\n")
}

// save connection string in safe place
// documentation, arrange, clear
// structure of project
// delete testcases function
// is all params come from user in update or not
// answer is need?

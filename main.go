package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"leetcoder/handlers"
)

func main() {

	router := gin.Default()

	// set url for CRUD
	router.GET("/", handlers.GetAllQuestions)
	router.GET("/questions", handlers.GetAllQuestions)
	router.GET("/questions/:id", handlers.GetQuestionById)
	router.POST("/questions", handlers.CreateQuestion)
	router.PUT("/questions/:id", handlers.UpdateQuestion)
	router.DELETE("/questions/:id", handlers.DeleteQuestion)
	router.POST("/questions/check-answer/:id", handlers.CheckAnswer)

	router.Run(":8080")

	fmt.Printf("server runs on localhost:8080\n")
}

// save connection string in safe place
// documentation, arrange, clear
// structure of project
// delete testcases function
// is all params come from user in update or not
// answer is need?

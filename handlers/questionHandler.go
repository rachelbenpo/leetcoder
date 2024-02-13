package handlers

import (
	"fmt"
	"leetcoder/models"
	"leetcoder/services"

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

	err := services.DeleteQuestion(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	c.JSON(200, gin.H{"message": "question deleted"})
}

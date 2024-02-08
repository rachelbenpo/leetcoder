package main

import (
	//"app/models"
	//"app/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func checkAnswer(c *gin.Context) {

	// get answer params
	var ans Answer

	err := c.BindJSON(&ans)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		fmt.Print(err)
		return
	}

	// get the  question
	questionID := c.Param("id")
	q, err := getQuestionById2(questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		fmt.Print(err)
		return
	}

	// run code
	result, err := checkAnswer2(q, ans)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		fmt.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"correct": result,
	})
}

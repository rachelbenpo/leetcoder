package handlers

import (
	//"app/models"
	//"app/utils"
	"fmt"
	"leetcoder/models"
	"leetcoder/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAnswer(c *gin.Context) {

	// get answer params
	var ans models.Answer

	err := c.BindJSON(&ans)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		fmt.Print(err)
		return
	}

	// get the  question
	questionID := c.Param("id")
	q, err := services.GetQuestionById(questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		fmt.Print(err)
		return
	}

	// run code
	result, err := services.CheckAnswer(ans, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		fmt.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"correct": result,
	})
}

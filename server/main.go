package main

import (
	"bufio"
	"fmt"
	"leetcode-server/DB"
	"leetcode-server/config"
	"leetcode-server/handlers"
	"leetcode-server/services"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	// get github credentials from user
	setConfigurations()

	err := DB.InitializeDB()
	if err != nil {
		fmt.Println(err)
	}

	// init image first so the user will make it public
	err = services.InitImage()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("please go to http://ghcr.io/" + config.UserName + "/checking-container/ and make the image public\n")

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

// set configurations
func setConfigurations() {

	input := bufio.NewScanner(os.Stdin)

	userName, isexist := os.LookupEnv("githubUsername")
	if !isexist {
		fmt.Println("Enter Your Github User Name: ")
		input.Scan()
		userName = input.Text()
		os.Setenv("githubUsername", userName)
	}

	token, isexist := os.LookupEnv("githubToken")
	if !isexist {
		fmt.Println("Enter Your Github Token: ")
		input.Scan()
		token = input.Text()
		os.Setenv("githubToken", token)
	}
	config.UserName = userName
	config.Token = token
}

package main

import "fmt"

func main() {

	var err error
	// run the CLI
	for {
		switch operation {

		case "Choose operation":
			if err = chooseOperation(); err != nil {
				fmt.Println("Error:", err)
				return
			}

		case "Get All Questions":
			if err = runGetAllForm(); err != nil {
				fmt.Println("error getting all questions from server: ", err)
				operation = "Choose operation"
			}

		case "Get Question by ID":
			if err = RunGetByIdForm(); err != nil {
				fmt.Println("error getting question from server: ", err)
				operation = "Choose operation"
			}

		case "Show single question":
			if err = runShowSingleForm(); err != nil {
				fmt.Println("error getting question from server: ", err)
				operation = "Choose operation"
			}

		case "Create Question":
			if err = runCreateQuestionForm(); err != nil {
				fmt.Println("error creating question: ", err)
			}
			operation = "Choose operation"

		case "Update Question":
			if err := runUpdateForm(); err != nil {
				fmt.Println("error updating question: ", err)
			}
			operation = "Choose operation"

		case "Delete Question":
			if err = runDeleteForm(); err != nil {
				fmt.Println("error deleting question:", err)
			}
			operation = "Choose operation"

		case "Check Answer":
			if err = runCheckForm(); err != nil {
				operation = "Choose operation"
				fmt.Println("error testing your answer: ", err)
			}
		case "Exit":
			return

		default:
			fmt.Println("Invalid operation selected.")
			return
		}
	}
}

package main

import (
	"database/sql"
	"fmt"
	"leetcode-server/config"

	_ "github.com/go-sql-driver/mysql"
)

// checks if the database exists, creates it if not, and returns a database connection.
func InitializeDB() (*sql.DB, error) {
    
	db, err := sql.Open("mysql", config.DBConnectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Check if the database exists
	err = db.Ping()
	if err == nil {
		fmt.Println("Database already exists.")
		return db, nil
	}

	// If the database doesn't exist, create it
	fmt.Println("Database does not exist. Creating...")
	err = createDatabase()
	if err != nil {
		return nil, err
	}

	// Retry connecting to the newly created database
	db, err = sql.Open("mysql", config.DBConnectionString)
	if err != nil {
		return nil, err
	}

	fmt.Println("Database created successfully.")
	return db, nil
}

func createDatabase() error {

	db, err := sql.Open("mysql", config.DBConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the database if not exists
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS leetcode")
	if err != nil {
		return err
	}

	// Connect to the leetcode database
	db, err = sql.Open("mysql", config.DBConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the 'questions' table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS leetcode.questions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name NVARCHAR(100) NOT NULL,
			instructions NVARCHAR(1000) NOT NULL,
			answer NVARCHAR(2000) NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create the 'test_cases' table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS leetcode.test_cases (
			id INT AUTO_INCREMENT PRIMARY KEY,
			question_id INT NOT NULL,
			input NVARCHAR(200) NULL,
			output NVARCHAR(200) NULL,
			CONSTRAINT fk_question_id FOREIGN KEY (question_id) REFERENCES leetcode.questions (id) ON DELETE NO ACTION ON UPDATE NO ACTION
		)
	`)
	if err != nil {
		return err
	}

	return nil
}
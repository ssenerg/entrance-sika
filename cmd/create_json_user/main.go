package main

import (
	"entrance/database"
	"log"
	"fmt"
)

func main() {

	dbConfig, err := database.NewConfig("localhost", "user", "password", "dbname", "5432", "disable", "UTC", nil)
	if err != nil {
		log.Fatal(fmt.Printf("Error creating db config: %v", err))
	}
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatal(fmt.Printf("Error initializing database: %v", err))
	}

	users, err := database.ReadUserFromJson("users_data.json")
	if err != nil {
		log.Fatal(fmt.Printf("Error reading users from json: %v", err))
	}
	db.CreateUsers(users)
}

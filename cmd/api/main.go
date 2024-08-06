package main

import (
	"entrance/api"
	"entrance/database"
	"log"
	"fmt"
)

func main() {

	dbConfig, err := database.NewConfig("localhost", "user", "password", "dbname", "5432", "disable", "UTC", nil)
	if err != nil {
		log.Fatal(fmt.Printf("Error creating db config: %v", err))
	}

	server, err := api.NewServer(":8080", dbConfig)
	if err != nil {
		log.Fatal(fmt.Printf("Error creating server: %v", err))
	}

	server.RegisterHandlers()
	server.Start()
}

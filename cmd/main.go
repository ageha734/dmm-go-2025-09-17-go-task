package main

import (
	"log"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/database"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/routes"
)

func main() {
	database.InitDatabase()

	router := routes.SetupRouter()

	log.Println("Starting server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

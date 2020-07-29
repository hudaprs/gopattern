package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopattern/config"
	"gopattern/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	ApiRoutes := routes.Api{}

	// Check env data
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load env")
		panic(err)
	}

	// Init the database
	config.Connect(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)

	// Init the routes
	ApiRoutes.ServeRoutes()

	// Run the server
	fmt.Println("Connected To Database")
	fmt.Println("Server started port 8000")
	log.Fatal(http.ListenAndServe(":8000", ApiRoutes.Router))
}

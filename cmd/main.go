// Package main is the entry point of our application
// In Go, the main package is special - it's where the program starts
package main

import (
	"log" // Standard library for logging

	"github.com/Asif-Faizal/Gommerce/cmd/api"
)

// main is the entry point function that gets called when the program starts
// It initializes and starts the API server
func main() {
	// Create a new API server instance
	// ":3000" means the server will listen on port 3000
	// nil is passed as the database connection (to be implemented later)
	server := api.NewAPIServer(":3000", nil)

	// Start the server and handle any potential errors
	if err := server.Run(); err != nil {
		// If there's an error, log it and exit the program
		log.Fatal(err)
	}
}

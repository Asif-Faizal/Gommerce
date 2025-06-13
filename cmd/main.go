// Package main is the entry point of our application
// In Go, the main package is special - it's where the program starts
package main

import (
	"log" // Standard library for logging

	"github.com/Asif-Faizal/Gommerce/cmd/api"
	"github.com/Asif-Faizal/Gommerce/config"
	"github.com/Asif-Faizal/Gommerce/db"
	"github.com/go-sql-driver/mysql"
)

// main is the entry point function that gets called when the program starts
// It initializes the database connection and starts the API server
func main() {
	// Initialize MySQL database connection using environment configuration
	db, err := db.MySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Net:                  "tcp",
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	// Create a new API server instance with the configured port and database connection
	server := api.NewAPIServer(config.Envs.Port, db)

	// Start the server and handle any potential errors
	if err := server.Run(); err != nil {
		// If there's an error, log it and exit the program
		log.Fatal(err)
	}
}

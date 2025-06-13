// Package main is the entry point of our application
// In Go, the main package is special - it's where the program starts
package main

import (
	"database/sql"
	"log" // Standard library for logging

	"github.com/Asif-Faizal/Gommerce/cmd/api"
	"github.com/Asif-Faizal/Gommerce/config"
	"github.com/Asif-Faizal/Gommerce/db"
	"github.com/go-sql-driver/mysql"
)

// main is the entry point function that gets called when the program starts
// It initializes the database connection and starts the API server
func main() {
	// Log the configuration being used
	log.Printf("Starting server with configuration:")
	log.Printf("Host: %s", config.Envs.PublicHost)
	log.Printf("Port: %s", config.Envs.Port)
	log.Printf("Database: %s@%s/%s", config.Envs.DBUser, config.Envs.DBAddress, config.Envs.DBName)

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

	initStorage(db)

	// Create a new API server instance with the configured port and database connection
	server := api.NewAPIServer(config.Envs.Port, db)

	// Start the server and handle any potential errors
	log.Printf("Starting server on %s", config.Envs.Port)
	if err := server.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func initStorage(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	log.Println("Successfully connected to database")
	return nil
}

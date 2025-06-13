// Package config handles all environment and configuration settings
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration values for the application
// These values can be set through environment variables or will use defaults
type Config struct {
	PublicHost string // The public host URL for the API
	Port       string // The port number the server will listen on
	DBUser     string // Database username
	DBPassword string // Database password
	DBAddress  string // Database host address and port
	DBName     string // Database name
}

// Envs is a global variable that holds the application configuration
// It's initialized when the package is imported
var Envs = InitConfig()

// InitConfig initializes the configuration with environment variables or default values
// Returns a Config struct with all necessary settings
func InitConfig() Config {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       ":" + getEnv("PORT", "8080"), // Add colon prefix for proper port format
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "root"),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:     getEnv("DB_NAME", "gommerce"),
	}
}

// getEnv retrieves an environment variable or returns a default value
// key: The name of the environment variable to look for
// defaultValue: The value to return if the environment variable is not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Package db handles all database-related functionality
package db

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

// MySQLStorage creates and returns a new MySQL database connection
// It takes a mysql.Config struct containing all necessary connection parameters
// Returns a *sql.DB connection and any potential error
func MySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	// Open a new database connection using the provided configuration
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}

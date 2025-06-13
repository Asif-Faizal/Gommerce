package main

import (
	"log"
	"os"

	"github.com/Asif-Faizal/Gommerce/config"
	"github.com/Asif-Faizal/Gommerce/db"
	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log.Printf("Starting server with configuration:")
	log.Printf("Host: %s", config.Envs.PublicHost)
	log.Printf("Port: %s", config.Envs.Port)
	log.Printf("Database: %s@%s/%s", config.Envs.DBUser, config.Envs.DBAddress, config.Envs.DBName)

	// Initialize MySQL database connection using environment configuration
	db, err := db.MySQLStorage(mysqldriver.Config{
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

	// Create MySQL migration driver
	driver, err := mysqlmigrate.WithInstance(db, &mysqlmigrate.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Create new migration instance
	m, err := migrate.NewWithDatabaseInstance("file://cmd/migrate/migrations", "mysql", driver)
	if err != nil {
		log.Fatal(err)
	}

	cdm := os.Args[(len(os.Args) - 1)]
	if cdm == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	} else if cdm == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}

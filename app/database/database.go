package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %s\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %s\n", err)
	}

	fmt.Println("Successfully connected to the database")
	createTablesIfNotExist()
}

func GetDB() *sql.DB {
	return db
}

func createTablesIfNotExist() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS buckets (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatalf("Error creating buckets table: %s\n", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS objects (
			id SERIAL PRIMARY KEY,
			bucket_id INTEGER REFERENCES buckets(id),
			name TEXT NOT NULL,
			data BYTEA,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatalf("Error creating objects table: %s\n", err)
	}

	fmt.Println("Tables are set up successfully")
}

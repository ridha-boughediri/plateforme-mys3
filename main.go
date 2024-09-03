package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found or could not load")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)
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

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func main() {
	initDB()

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

// main.go
package main

import (
	"fmt"
	"net/http"
	"os"

	"example.com/hello/app/database"
	"example.com/hello/app/storage"
)

func main() {
	database.InitDB()

	http.HandleFunc("/", storage.GetRoot)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/hello/app/database"
	"example.com/hello/app/storage"
	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: Method=%s, URL=%s, Path=%s, RemoteAddr=%s, UserAgent=%s",
			r.Method, r.URL.String(), r.URL.Path, r.RemoteAddr, r.UserAgent())

		for name, values := range r.Header {
			log.Printf("Header: %s = %s", name, values)
		}

		next.ServeHTTP(w, r)

		log.Printf("Handled request: Method=%s, URL=%s", r.Method, r.URL.String())
	})
}

func main() {
	database.InitDB()

	router := mux.NewRouter()

	// Apply the logging middleware to all routes
	router.Use(loggingMiddleware)

	// Bucket routes
	router.HandleFunc("/", storage.ListBuckets).Methods(http.MethodGet)
	router.HandleFunc("/{bucketName}/", storage.CreateBucket).Methods(http.MethodPut)
	router.HandleFunc("/{bucketName

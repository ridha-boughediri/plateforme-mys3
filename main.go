package main

import (
	"log"
	"net/http"
	"example.com/hello/app/config"
	"example.com/hello/app/controller"
	"example.com/hello/app/database"
	"example.com/hello/app/middleware"
	"github.com/gorilla/mux"
)

func main() {
	// Load environment variables
	config.LoadConfig()
	
	// Initialize the database
	database.InitDB()

	// Set up the router
	router := mux.NewRouter()

	// Apply logging middleware
	router.Use(middleware.LoggingMiddleware)

	// Bucket routes
	router.HandleFunc("/", controller.ListBuckets).Methods(http.MethodGet)
	router.HandleFunc("/{bucketName}/", controller.CreateBucket).Methods(http.MethodPut)
	router.HandleFunc("/{bucketName}/", controller.DeleteBucket).Methods(http.MethodDelete)

	// Object routes
	router.HandleFunc("/{bucketName}/objects", controller.ListObjects).Methods(http.MethodGet)
	router.HandleFunc("/{bucketName}/objects/{objectName}", controller.AddObject).Methods(http.MethodPut)
	router.HandleFunc("/{bucketName}/objects/{objectName}", controller.DownloadObject).Methods(http.MethodGet)
	router.HandleFunc("/{bucketName}/objects/{objectName}", controller.DeleteObject).Methods(http.MethodDelete)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}

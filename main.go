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
	config.LoadConfig()
	
	database.InitDB()

	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)

	router.HandleFunc("/", controller.ListBuckets).Methods(http.MethodGet)
	router.HandleFunc("/{bucketName}/", controller.CreateBucket).Methods(http.MethodPut)
	router.HandleFunc("/{bucketName}/", controller.DeleteBucket).Methods(http.MethodDelete)

	router.HandleFunc("/{bucketName}/objects", controller.ListObjects).Methods(http.MethodGet)
	router.HandleFunc("/{bucketName}/objects/{objectName}", controller.AddObject).Methods(http.MethodPut)
	router.HandleFunc("/{bucketName}/objects/{objectName}", controller.DownloadObject).Methods(http.MethodGet)
	router.HandleFunc("/{bucketName}/objects/{objectName}", controller.DeleteObject).Methods(http.MethodDelete)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}

package main

import (
	"log"
	"net/http"
	"plateforme-mys3/handlers"
	"plateforme-mys3/storage"

	"github.com/gorilla/mux"
)

func main() {
	db := storage.InitDB("storage.db")
	storage.Migrate(db)

	r := mux.NewRouter()

	r.HandleFunc("/products", handlers.GetProducts(db)).Methods("GET")
	r.HandleFunc("/products", handlers.CreateProduct(db)).Methods("POST")
	r.HandleFunc("/products/{id}", handlers.UpdateProduct(db)).Methods("PUT")
	r.HandleFunc("/products/{id}", handlers.DeleteProduct(db)).Methods("DELETE")

	log.Println("API is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

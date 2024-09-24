// cmd/main.go
package main

import (
	"log"
	"net/http"
	"plateforme-mys3/config"
	"plateforme-mys3/internal/handlers"
	"plateforme-mys3/internal/middleware"
	"plateforme-mys3/internal/storage"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	s := storage.NewStorage(cfg.StoragePath)

	r := mux.NewRouter()

	// Routes pour les buckets
	r.HandleFunc("/", handlers.ListBucketsHandler(s)).Methods("GET")
	r.HandleFunc("/{bucket}", handlers.BucketHandler(s)).Methods("PUT", "DELETE")

	// Routes pour les objets
	r.HandleFunc("/{bucket}/{object}", handlers.ObjectHandler(s)).Methods("PUT", "GET", "DELETE")
	r.HandleFunc("/{bucket}/list-objects", handlers.ListObjectsHandler(s)).Methods("GET")

	// Appliquer le middleware d'authentification
	r.Use(func(next http.Handler) http.Handler {
		return middleware.AuthMiddleware(next, cfg)
	})

	log.Println("Serveur démarré sur le port 9000")
	if err := http.ListenAndServe(":9000", r); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur : %v", err)
	}
}

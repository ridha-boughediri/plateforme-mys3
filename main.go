package main

import (
	"log"
	"os"
	"plateforme-mys3/config"
	"plateforme-mys3/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Initialiser le client MinIO
	err := config.InitMinioClient()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation du client MinIO: %v", err)
	}

	// Définir les routes pour les handlers
	app.Post("/upload-file", handlers.UploadFile)
	app.Get("/create-bucket", handlers.CreateBucketHandler)
	app.Get("/list-buckets", handlers.ListBucketsHandler)
	app.Delete("/delete-bucket", handlers.DeleteBucketHandler)

	// Démarrer le serveur
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Port par défaut
	}
	log.Printf("Serveur démarré sur le port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}

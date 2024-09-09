package main

import (
	"log"
	"os"

	"plateforme-mys3/config"
	"plateforme-mys3/controllers"

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
	app.Post("/upload-file", controllers.UploadFile)
	app.Get("/create-bucket", controllers.CreateBucketHandler)
	app.Get("/delete-file", controllers.DeleteFileHandler)
	app.Get("/download-file", controllers.DownloadFileHandler)

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

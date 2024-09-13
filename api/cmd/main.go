package main

import (
	"plateforme-mys3/api/handlers"

	"github.com/gofiber/fiber/v2"
)

const storagePath = "./storage"

func main() {
	app := fiber.New()

	// Endpoint pour la vérification de l'état du serveur
	app.All("/probe-bsign*", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/xml")
		return c.Status(fiber.StatusOK).SendString("<Response></Response>")
	})

	// Route pour la création d'un bucket
	app.Put("/bucket/:bucketName", handlers.CreateBucketHandler)

	// Route pour lister les buckets
	app.Get("/buckets", handlers.ListBucketsHandler)

	// Route pour uploader un objet dans un bucket
	app.Put("/bucket/:bucketName/:objectName", handlers.UploadObjectHandler)

	// Route pour télécharger un objet d'un bucket
	app.Get("/bucket/:bucketName/:objectName", handlers.DownloadObjectHandler)

	// Route pour supprimer un objet d'un bucket
	app.Delete("/bucket/:bucketName/:objectName", handlers.DeleteObjectHandler)

	// Route pour lister les objets dans un bucket
	app.Get("/bucket/:bucketName/objects", handlers.ListObjectsHandler)

	// Démarrer le serveur sur le port 3000
	if err := app.Listen(":9000"); err != nil {
		panic("Erreur lors du démarrage du serveur : " + err.Error())
	}
}

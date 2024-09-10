package controllers

import (
	"context"
	"os"

	"plateforme-mys3/config"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func DeleteFileHandler(c *fiber.Ctx) error {
	// Récupérer le nom du fichier depuis les paramètres de requête
	fileName := c.Query("file")
	bucketName := os.Getenv("MINIO_BUCKET")

	// Vérifier si le bucket a un nom par défaut
	if bucketName == "" {
		bucketName = "default-bucket" // Utiliser un nom par défaut ou gérer une absence de bucket
	}

	// Vérification si le nom du fichier est fourni
	if fileName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Le nom du fichier est requis",
		})
	}

	// Créer un contexte pour l'opération MinIO
	ctx := context.Background()

	// Vérification si le client MinIO est initialisé correctement
	minioClient := config.MinioClient
	if minioClient == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Le client MinIO n'est pas disponible",
		})
	}

	// Suppression de l'objet (fichier)
	err := minioClient.RemoveObject(ctx, bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Erreur lors de la suppression du fichier : " + err.Error(),
		})
	}

	// Réponse en cas de succès
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "Fichier supprimé avec succès",
	})
}

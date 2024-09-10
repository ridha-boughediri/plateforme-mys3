package controllers

import (
	"context"
	"os"

	"plateforme-mys3/config"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func DownloadFileHandler(c *fiber.Ctx) error {
	fileName := c.Query("file")
	bucketName := os.Getenv("MINIO_BUCKET")
	if bucketName == "" {
		bucketName = "default-bucket" // Utilisez un nom par défaut ou gérez le cas où le bucket est manquant
	}

	if fileName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Le nom du fichier est requis",
		})
	}

	ctx := context.Background()
	minioClient := config.MinioClient

	object, err := minioClient.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Erreur lors de la récupération du fichier",
		})
	}
	defer object.Close()

	// Renvoyez le contenu du fichier en streaming
	return c.SendStream(object)
}

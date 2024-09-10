package controllers

import (
	"context"
	"plateforme-mys3/config"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

// CreateBucketHandler gère la création de buckets
func CreateBucketHandler(c *fiber.Ctx) error {
	bucketName := c.Query("bucket")

	// Vérifiez si le client Minio est initialisé
	if config.MinioClient == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Client MinIO non initialisé",
		})
	}

	ctx := context.Background()

	// Appel à la méthode MakeBucket via l'interface
	err := config.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Erreur lors de la création du bucket",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   "Bucket créé avec succès",
	})
}

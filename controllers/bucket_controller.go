package controllers

import (
	"context"
	"plateforme-mys3/config"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func CreateBucketHandler(c *fiber.Ctx) error {
	bucketName := c.Query("bucket")

	// Utiliser le client MinIO global défini dans le package config
	err := config.MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "Erreur lors de la création du bucket",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg": "Bucket créé avec succès",
	})
}

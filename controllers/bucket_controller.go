package controllers

import (
	"context"
	"log"
	"plateforme-mys3/config"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func CreateBucketHandler(c *fiber.Ctx) error {
	bucketName := c.Query("bucket")

	if bucketName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Le nom du bucket est requis",
		})
	}

	ctx := context.Background()
	if config.MinioClient == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Client MinIO non initialisé",
		})
	}

	err := config.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	log.Printf("Bucket %s créé avec succès\n", bucketName)
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Bucket créé avec succès",
	})
}

package controllers

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

// DeleteFileHandler supprime un fichier de MinIO
func DeleteFileHandler(c *fiber.Ctx) error {
	fileName := c.Query("file")
	bucketName := os.Getenv("MINIO_BUCKET")

	if fileName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Le nom du fichier est requis",
		})
	}

	ctx := context.Background()
	minioClient, err := initMinioClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	err = minioClient.RemoveObject(ctx, bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "Fichier supprimé avec succès",
	})
}

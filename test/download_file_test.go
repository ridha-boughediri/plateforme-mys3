package controllers

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

// DownloadFileHandler télécharge un fichier de MinIO
func DownloadFileHandler(c *fiber.Ctx) error {
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

	object, err := minioClient.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	defer object.Close()

	return c.SendStream(object)
}

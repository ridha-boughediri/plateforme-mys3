package handlers

import (
	"context"
	"log"
	"os"

	"plateforme-mys3/config"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func UploadFile(c *fiber.Ctx) error {
	ctx := context.Background()
	bucketName := os.Getenv("MINIO_BUCKET")
	if bucketName == "" {
		bucketName = "default-bucket" // Utilisez un nom par défaut ou gérez le cas où le bucket est manquant
	}

	file, err := c.FormFile("fileUpload")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Fichier non trouvé dans la requête",
		})
	}

	buffer, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Erreur lors de l'ouverture du fichier",
		})
	}
	defer buffer.Close()

	// Utilisez le client MinIO centralisé
	minioClient := config.MinioClient

	objectName := file.Filename
	fileBuffer := buffer
	contentType := file.Header["Content-Type"][0]
	fileSize := file.Size

	// Upload the file with PutObject
	info, err := minioClient.PutObject(ctx, bucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Erreur lors de l'upload du fichier",
		})
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Fichier uploadé avec succès",
		"info":  info,
	})
}

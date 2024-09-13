package handlers

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"plateforme-mys3/api/dto"
	"plateforme-mys3/api/services" // Import pour utiliser StorageBasePath

	"github.com/gofiber/fiber/v2"
)

// UploadObjectHandler permet d'uploader un fichier dans un bucket
func UploadObjectHandler(c *fiber.Ctx) error {
	bucketName := c.Params("bucketName")
	objectName := c.Params("objectName")
	bucketPath := filepath.Join(services.StorageBasePath, bucketName)

	// Vérifier si le bucket existe
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).SendString("Le bucket n'existe pas")
	}

	// Enregistrer le fichier uploadé
	filePath := filepath.Join(bucketPath, objectName)
	file, err := os.Create(filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Erreur lors de l'upload du fichier")
	}
	defer file.Close()

	// Convertir []byte en io.Reader
	buffer := bytes.NewReader(c.Body())

	_, err = io.Copy(file, buffer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Erreur lors de l'upload")
	}

	return c.Status(fiber.StatusOK).SendString("Fichier uploadé avec succès")
}

// Dans le fichier handlers/objects.go

// DownloadObjectHandler permet de télécharger un fichier depuis un bucket
func DownloadObjectHandler(c *fiber.Ctx) error {
	bucketName := c.Params("bucketName")
	objectName := c.Params("objectName")
	filePath := filepath.Join(services.StorageBasePath, bucketName, objectName)

	// Vérifier si le fichier existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).SendString("Fichier non trouvé")
	}

	return c.SendFile(filePath)
}

// DeleteObjectHandler permet de supprimer un fichier d'un bucket
func DeleteObjectHandler(c *fiber.Ctx) error {
	bucketName := c.Params("bucketName")
	objectName := c.Params("objectName")
	filePath := filepath.Join(services.StorageBasePath, bucketName, objectName)

	// Vérifier si le fichier existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).SendString("Fichier non trouvé")
	}

	// Supprimer le fichier
	if err := os.Remove(filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Erreur lors de la suppression du fichier")
	}

	return c.Status(fiber.StatusOK).SendString("Fichier supprimé avec succès")
}

// ListObjectsHandler permet de lister tous les objets dans un bucket
func ListObjectsHandler(c *fiber.Ctx) error {
	bucketName := c.Params("bucketName")
	bucketPath := filepath.Join(services.StorageBasePath, bucketName)

	// Vérifier si le bucket existe
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		errorResponse := dto.ErrorResponse{
			Code:      "NoSuchBucket",
			Message:   "The specified bucket does not exist.",
			Resource:  bucketName,
			RequestID: "12345",
		}
		xmlResponse, _ := xml.MarshalIndent(errorResponse, "", "  ")
		c.Set("Content-Type", "application/xml")
		return c.Status(fiber.StatusNotFound).Send(xmlResponse)
	}

	// Lister les objets dans le bucket
	files, err := os.ReadDir(bucketPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Erreur lors de la récupération des objets")
	}

	// Construire la réponse XML
	var objects []dto.ObjectInfo
	for _, file := range files {
		if !file.IsDir() {
			fileInfo, _ := file.Info()
			object := dto.ObjectInfo{
				Key:          file.Name(),
				LastModified: fileInfo.ModTime().Format("2006-01-02T15:04:05Z"),
				Size:         fileInfo.Size(),
			}
			objects = append(objects, object)
		}
	}

	response := dto.ListObjectsResponse{
		Name:        bucketName,
		Prefix:      "",
		Marker:      "",
		MaxKeys:     len(objects),
		IsTruncated: false,
		Contents:    objects,
	}

	// Convertir en XML
	xmlResponse, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Erreur lors de la génération de la réponse XML")
	}

	c.Set("Content-Type", "application/xml")
	return c.Send(xmlResponse)
}

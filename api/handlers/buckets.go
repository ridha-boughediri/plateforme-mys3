package handlers

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"plateforme-mys3/api/dto"
	"plateforme-mys3/api/services"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

// Validation des noms de bucket
var validBucketName = regexp.MustCompile(`^[a-z0-9.\-]{3,63}$`)

// CreateBucketHandler crée un nouveau bucket en respectant les conventions S3
func CreateBucketHandler(c *fiber.Ctx) error {
	bucketName := c.Params("bucketName")

	// Valider le nom du bucket
	if !validBucketName.MatchString(bucketName) {
		errorResponse := dto.ErrorResponse{
			Code:      "InvalidBucketName",
			Message:   "The specified bucket is not valid.",
			Resource:  bucketName,
			RequestID: "12345",
		}
		xmlResponse, _ := xml.MarshalIndent(errorResponse, "", "  ")
		c.Set("Content-Type", "application/xml")
		return c.Status(fiber.StatusBadRequest).Send(xmlResponse)
	}

	bucketPath := filepath.Join(services.StorageBasePath, bucketName)

	// Vérifier si le bucket existe déjà
	if _, err := os.Stat(bucketPath); !os.IsNotExist(err) {
		errorResponse := dto.ErrorResponse{
			Code:      "BucketAlreadyExists",
			Message:   "The requested bucket name is not available.",
			Resource:  bucketName,
			RequestID: "12345",
		}
		xmlResponse, _ := xml.MarshalIndent(errorResponse, "", "  ")
		c.Set("Content-Type", "application/xml")
		return c.Status(fiber.StatusConflict).Send(xmlResponse)
	}

	// Créer le bucket
	if err := os.MkdirAll(bucketPath, 0755); err != nil {
		errorResponse := dto.ErrorResponse{
			Code:      "InternalError",
			Message:   "We encountered an internal error. Please try again.",
			Resource:  bucketName,
			RequestID: "12345",
		}
		xmlResponse, _ := xml.MarshalIndent(errorResponse, "", "  ")
		c.Set("Content-Type", "application/xml")
		return c.Status(fiber.StatusInternalServerError).Send(xmlResponse)
	}

	// Réponse de succès
	return c.Status(fiber.StatusOK).SendString("Bucket créé avec succès")
}

// ListBucketsHandler permet de lister tous les buckets

func ListBucketsHandler(c *fiber.Ctx) error {
	// Ouvre le répertoire où les buckets sont stockés
	entries, err := os.ReadDir(services.StorageBasePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Erreur lors de la récupération des buckets")
	}

	// Créer la réponse avec les buckets
	var bucketList []dto.BucketInfo
	for _, entry := range entries {
		if entry.IsDir() {
			info, err := entry.Info() // Correction ici
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Erreur lors de la récupération des informations du bucket")
			}

			bucketInfo := dto.BucketInfo{
				Name:         entry.Name(),
				CreationDate: info.ModTime().Format("2006-01-02T15:04:05Z"),
			}
			bucketList = append(bucketList, bucketInfo)
		}
	}

	// Construire la réponse XML
	response := dto.ListBucketsResponse{
		Owner: dto.Owner{
			ID:          "12345", // ID fictif du propriétaire
			DisplayName: "mon-utilisateur",
		},
		Buckets: bucketList,
	}

	// Convertir en XML
	xmlResponse, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Erreur lors de la génération de la réponse XML")
	}

	c.Set("Content-Type", "application/xml")
	return c.Send(xmlResponse)
}

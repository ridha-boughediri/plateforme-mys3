package handlers

import (
	"context"
	"log"
	"plateforme-mys3/config"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

// CreateBucketHandler permet de créer un nouveau bucket
func CreateBucketHandler(c *fiber.Ctx) error {
	bucketName := c.Query("bucket")

	err := config.MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		log.Printf("Erreur lors de la création du bucket: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "Erreur lors de la création du bucket",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg": "Bucket créé avec succès",
	})
}

// ListBucketsHandler permet de lister tous les buckets
func ListBucketsHandler(c *fiber.Ctx) error {
	buckets, err := config.MinioClient.ListBuckets(context.Background())
	if err != nil {
		log.Printf("Erreur lors de la récupération des buckets: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "Erreur lors de la récupération des buckets",
		})
	}

	// Préparer une réponse avec tous les buckets
	var bucketList []string
	for _, bucket := range buckets {
		bucketList = append(bucketList, bucket.Name)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"buckets": bucketList,
	})
}

// DeleteBucketHandler permet de supprimer un bucket
func DeleteBucketHandler(c *fiber.Ctx) error {
	bucketName := c.Query("bucket")

	err := config.MinioClient.RemoveBucket(context.Background(), bucketName)
	if err != nil {
		log.Printf("Erreur lors de la suppression du bucket: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "Erreur lors de la suppression du bucket",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": "Bucket supprimé avec succès",
	})
}

// UpdateBucketHandler (facultatif dans MinIO) - Il n'y a pas de mise à jour directe des buckets dans MinIO
// Cependant, tu peux renommer un bucket ou modifier ses permissions. MinIO ne supporte pas directement le renommage,
// donc il faudra copier les objets vers un nouveau bucket et supprimer l'ancien. Si tu as besoin de cette fonctionnalité, il faudra
// la gérer manuellement en utilisant des méthodes comme `CopyObject` puis `RemoveBucket`.

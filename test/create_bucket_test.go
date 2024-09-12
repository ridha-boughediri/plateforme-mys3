package test

import (
	"context"
	"log"
	"plateforme-mys3/handlers"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestCreateBucketHandler_Success(t *testing.T) {
	// Initialiser le client Minio avec des informations d'accès
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("accessKey", "secretKey", ""),
		Secure: false, // false si vous n'utilisez pas SSL
	})
	if err != nil {
		t.Fatalf("Erreur lors de la création du client Minio : %v", err)
	}

	// Nom du bucket à créer
	bucketName := "test-bucket"

	// Appel au handler réel avec le client réel
	err = handlers.CreateBucketHandler(minioClient, bucketName)
	if err != nil {
		t.Errorf("Erreur lors de la création du bucket : %v", err)
	}

	// Nettoyage après le test (suppression du bucket)
	defer func() {
		err = minioClient.RemoveBucket(context.Background(), bucketName)
		if err != nil {
			log.Fatalf("Erreur lors de la suppression du bucket : %v", err)
		}
	}()
}

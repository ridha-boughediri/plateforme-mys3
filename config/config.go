package config

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Global MinIO client
var MinioClient *minio.Client

// Initialize MinIO Client
func InitMinioClient() error {
	endpoint := "localhost:9000"
	accessKeyID := "admin"
	secretAccessKey := "admin1234@@@test"
	useSSL := false

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	MinioClient = client
	return nil
}

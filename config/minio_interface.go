package config

import (
	"context"

	"github.com/minio/minio-go/v7"
)

// MinioClientInterface représente l'interface des méthodes utilisées par MinIO
type MinioClientInterface interface {
	MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error
}

// MinioClient est l'implémentation réelle de l'interface
var MinioClient MinioClientInterface

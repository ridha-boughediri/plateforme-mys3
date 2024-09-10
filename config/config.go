package config

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClientWrapper est un adaptateur pour le vrai client MinIO
type MinioClientWrapper struct {
	Client *minio.Client
}

var MinioClient MinioClientInterface

// Implémentation des méthodes de l'interface MinioClientInterface

func (m *MinioClientWrapper) MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error {
	return m.Client.MakeBucket(ctx, bucketName, opts)
}

func (m *MinioClientWrapper) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	return m.Client.BucketExists(ctx, bucketName)
}

func (m *MinioClientWrapper) RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	return m.Client.RemoveObject(ctx, bucketName, objectName, opts)
}

func (m *MinioClientWrapper) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	return m.Client.GetObject(ctx, bucketName, objectName, opts)
}

func (m *MinioClientWrapper) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	return m.Client.PutObject(ctx, bucketName, objectName, reader, objectSize, opts)
}

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

	// Créez un adaptateur pour le client MinIO réel
	MinioClient = &MinioClientWrapper{Client: client}
	return nil
}

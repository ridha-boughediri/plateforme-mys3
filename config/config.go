package config

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClientInterface définit les méthodes pour le client MinIO
type MinioClientInterface interface {
	MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error)
	ListBuckets(ctx context.Context) ([]minio.BucketInfo, error)
	RemoveBucket(ctx context.Context, bucketName string) error
}

// MinioClientWrapper est un adaptateur pour le client MinIO réel
type MinioClientWrapper struct {
	Client *minio.Client
}

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

func (m *MinioClientWrapper) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	return m.Client.ListBuckets(ctx)
}

func (m *MinioClientWrapper) RemoveBucket(ctx context.Context, bucketName string) error {
	return m.Client.RemoveBucket(ctx, bucketName)
}

// MinioClient est une variable globale qui stocke le client MinIO
var MinioClient MinioClientInterface

// InitMinioClient initialise le client MinIO
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

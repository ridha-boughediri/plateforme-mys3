import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestCreateBucketHandler_Success(t *testing.T) {
	// Créer un vrai client Minio
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("accessKey", "secretKey", ""),
		Secure: false,
	})
	if err != nil {
		t.Fatalf("failed to create minio client: %v", err)
	}

	// Appeler le handler avec le vrai client Minio
	err = handlers.CreateBucketHandler(minioClient, "my-test-bucket")
	if err != nil {
		t.Errorf("CreateBucketHandler failed: %v", err)
	}

	// Nettoyage : supprimer le bucket après le test
	defer minioClient.RemoveBucket(context.Background(), "my-test-bucket")
}

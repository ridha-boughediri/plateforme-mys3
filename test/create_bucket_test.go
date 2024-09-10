package test

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"plateforme-mys3/config"
	"plateforme-mys3/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMinioClient est une structure qui simule le client MinIO
type MockMinioClient struct {
	mock.Mock
}

// Simule la méthode MakeBucket de MinIO
func (m *MockMinioClient) MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error {
	args := m.Called(ctx, bucketName, opts)
	return args.Error(0)
}

func TestCreateBucketHandler(t *testing.T) {
	app := fiber.New()

	// Définir le handler
	app.Post("/create-bucket", controllers.CreateBucketHandler)

	// Créer un mock pour MinIO
	mockMinioClient := new(MockMinioClient)
	// Remplacez le client MinIO réel par le mock
	config.MinioClient = mockMinioClient

	// Configurer le mock pour retourner un succès lors de la création d'un bucket
	mockMinioClient.On("MakeBucket", mock.Anything, "test-bucket", mock.Anything).Return(nil)

	// Créer la requête de test
	req, _ := http.NewRequest("POST", "/create-bucket?bucket=test-bucket", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode) // Vérifie que le code 201 est retourné

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Vérifier que la réponse contient le message attendu
	assert.Contains(t, string(body), `"msg":"Bucket créé avec succès"`)

	// Vérifier que le mock a bien été appelé
	mockMinioClient.AssertExpectations(t)
}

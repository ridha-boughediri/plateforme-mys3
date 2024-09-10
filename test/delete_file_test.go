package controllers

import (
	"context"
	"net/http"
	"os"
	"testing"

	"plateforme-mys3/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMinioClient est une structure qui simule le client MinIO
type MockMinioClient struct {
	mock.Mock
}

func (m *MockMinioClient) RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	args := m.Called(ctx, bucketName, objectName, opts)
	return args.Error(0)
}

func TestDeleteFileHandler(t *testing.T) {
	// Initialisation de Fiber
	app := fiber.New()

	// Définir le handler pour le test
	app.Delete("/delete", DeleteFileHandler)

	// Configuration de la variable d'environnement
	os.Setenv("MINIO_BUCKET", "test-bucket")

	// Mock MinIO client
	mockMinioClient := new(MockMinioClient)
	// Préparer le mock pour RemoveObject
	mockMinioClient.On("RemoveObject", mock.Anything, "test-bucket", "test-file.txt", minio.RemoveObjectOptions{}).Return(nil)

	// Remplacer le client MinIO réel par le mock
	config.MinioClient = mockMinioClient

	// Requête DELETE pour tester le handler
	req, err := http.NewRequest("DELETE", "/delete?file=test-file.txt", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Vérification de la réponse JSON
	body := utils.Body(resp)
	assert.Contains(t, string(body), `"error":false`)
	assert.Contains(t, string(body), `"msg":"Fichier supprimé avec succès"`)

	// Vérifier les appels mockés
	mockMinioClient.AssertExpectations(t)
}

func TestDeleteFileHandler_MissingFileName(t *testing.T) {
	// Initialisation de Fiber
	app := fiber.New()

	// Définir le handler pour le test
	app.Delete("/delete", DeleteFileHandler)

	// Configuration de la variable d'environnement
	os.Setenv("MINIO_BUCKET", "test-bucket")

	// Mock MinIO client
	mockMinioClient := new(MockMinioClient)
	config.MinioClient = mockMinioClient

	// Requête DELETE pour tester le handler sans nom de fichier
	req, err := http.NewRequest("DELETE", "/delete", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Vérification de la réponse JSON
	body := utils.Body(resp)
	assert.Contains(t, string(body), `"error":true`)
	assert.Contains(t, string(body), `"msg":"Le nom du fichier est requis"`)
}

func TestDeleteFileHandler_MissingMinioClient(t *testing.T) {
	// Initialisation de Fiber
	app := fiber.New()

	// Définir le handler pour le test
	app.Delete("/delete", DeleteFileHandler)

	// Configuration de la variable d'environnement
	os.Setenv("MINIO_BUCKET", "test-bucket")

	// Réinitialiser MinIO client pour simuler une absence
	config.MinioClient = nil

	// Requête DELETE pour tester le handler avec un client MinIO manquant
	req, err := http.NewRequest("DELETE", "/delete?file=test-file.txt", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	// Vérification de la réponse JSON
	body := utils.Body(resp)
	assert.Contains(t, string(body), `"error":true`)
	assert.Contains(t, string(body), `"msg":"Le client MinIO n'est pas disponible"`)
}

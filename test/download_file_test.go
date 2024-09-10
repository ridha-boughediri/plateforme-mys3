package controllers

import (
	"context"
	"io/ioutil"
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

func (m *MockMinioClient) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (minio.Object, error) {
	args := m.Called(ctx, bucketName, objectName, opts)
	return args.Get(0).(minio.Object), args.Error(1)
}

// MockObject est une structure qui simule l'objet MinIO retourné par le client
type MockObject struct {
	mock.Mock
}

func (m *MockObject) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockObject) Read(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func TestDownloadFileHandler(t *testing.T) {
	// Initialisation de Fiber
	app := fiber.New()

	// Définir le handler pour le test
	app.Get("/download", DownloadFileHandler)

	// Configuration de la variable d'environnement
	os.Setenv("MINIO_BUCKET", "test-bucket")

	// Mock MinIO client et objet
	mockMinioClient := new(MockMinioClient)
	mockObject := new(MockObject)
	mockObject.On("Read", mock.Anything).Return(len([]byte("file content")), nil)
	mockObject.On("Close").Return(nil)
	mockMinioClient.On("GetObject", mock.Anything, "test-bucket", "test-file.txt", minio.GetObjectOptions{}).Return(mockObject, nil)

	// Remplacer le client MinIO réel par le mock
	config.MinioClient = mockMinioClient

	// Requête GET pour tester le handler
	req, err := http.NewRequest("GET", "/download?file=test-file.txt", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Vérification du contenu de la réponse
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "file content", string(body))

	// Vérifier les appels mockés
	mockMinioClient.AssertExpectations(t)
	mockObject.AssertExpectations(t)
}

func TestDownloadFileHandler_MissingFileName(t *testing.T) {
	// Initialisation de Fiber
	app := fiber.New()

	// Définir le handler pour le test
	app.Get("/download", DownloadFileHandler)

	// Configuration de la variable d'environnement
	os.Setenv("MINIO_BUCKET", "test-bucket")

	// Mock MinIO client
	mockMinioClient := new(MockMinioClient)
	config.MinioClient = mockMinioClient

	// Requête GET pour tester le handler sans nom de fichier
	req, err := http.NewRequest("GET", "/download", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Vérification de la réponse JSON
	body := utils.Body(resp)
	assert.Contains(t, string(body), `"error":true`)
	assert.Contains(t, string(body), `"msg":"Le nom du fichier est requis"`)
}

func TestDownloadFileHandler_MissingMinioClient(t *testing.T) {
	// Initialisation de Fiber
	app := fiber.New()

	// Définir le handler pour le test
	app.Get("/download", DownloadFileHandler)

	// Configuration de la variable d'environnement
	os.Setenv("MINIO_BUCKET", "test-bucket")

	// Réinitialiser MinIO client pour simuler une absence
	config.MinioClient = nil

	// Requête GET pour tester le handler avec un client MinIO manquant
	req, err := http.NewRequest("GET", "/download?file=test-file.txt", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	// Vérification de la réponse JSON
	body := utils.Body(resp)
	assert.Contains(t, string(body), `"error":true`)
	assert.Contains(t, string(body), `"msg":"Le client MinIO n'est pas disponible"`)
}

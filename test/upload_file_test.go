package controllers

import (
	"bytes"
	"context"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"plateforme-mys3/config"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Interface qui représente les méthodes de MinIO utilisées dans votre code
type MinioClientInterface interface {
	PutObject(ctx context.Context, bucketName, objectName string, reader *bytes.Buffer, size int64, opts minio.PutObjectOptions) (minio.UploadInfo, error)
}

// MockMinioClient implémente l'interface MinioClientInterface
type MockMinioClient struct {
	mock.Mock
}

func (m *MockMinioClient) PutObject(ctx context.Context, bucketName, objectName string, reader *bytes.Buffer, size int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	args := m.Called(ctx, bucketName, objectName, reader, size, opts)
	return args.Get(0).(minio.UploadInfo), args.Error(1)
}

// TestUploadFile teste le handler UploadFile pour un upload réussi
func TestUploadFile(t *testing.T) {
	// Initialisation de Fiber
	app := fiber.New()

	// Définir le handler pour le test
	app.Post("/upload", UploadFile)

	// Configuration de la variable d'environnement
	os.Setenv("MINIO_BUCKET", "test-bucket")

	// Mock MinIO client
	mockMinioClient := new(MockMinioClient)
	config.MinioClient = mockMinioClient

	// Préparer le mock pour PutObject
	mockMinioClient.On("PutObject", mock.Anything, "test-bucket", "test-file.txt", mock.Anything, int64(14), minio.PutObjectOptions{ContentType: "text/plain"}).
		Return(minio.UploadInfo{Size: 14}, nil)

	// Créer un fichier de test pour l'upload
	fileContent := "test file content"
	fileBuffer := bytes.NewBufferString(fileContent)

	// Créer une requête POST pour uploader le fichier
	req, err := http.NewRequest("POST", "/upload", fileBuffer)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=---boundary")

	// Correction : Utiliser uniquement le nom du fichier pour `FormFile`
	// Suppression des arguments supplémentaires
	req.MultipartForm = &multipart.Form{
		File: map[string][]*multipart.FileHeader{
			"fileUpload": {
				&multipart.FileHeader{
					Filename: "test-file.txt",
					Size:     int64(len(fileContent)),
				},
			},
		},
	}

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Lire la réponse avec ioutil.ReadAll au lieu de utils.Body
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), `"error":false`)
	assert.Contains(t, string(body), `"msg":null`)

	// Vérifier les appels mockés
	mockMinioClient.AssertExpectations(t)
}

func TestUploadFile_MissingFile(t *testing.T) {
	// Initialisation de Fiber
	app := fiber.New()

	// Définir le handler pour le test
	app.Post("/upload", UploadFile)

	// Configuration de la variable d'environnement
	os.Setenv("MINIO_BUCKET", "test-bucket")

	// Mock MinIO client
	mockMinioClient := new(MockMinioClient)
	config.MinioClient = mockMinioClient

	// Requête POST pour tester le handler sans fichier
	req, err := http.NewRequest("POST", "/upload", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Lire la réponse avec ioutil.ReadAll au lieu de utils.Body
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), `"error":true`)
	assert.Contains(t, string(body), `"msg":"multipart: form-data boundary missing"`)
}

func TestUploadFile_MinioError(t *testing.T) {
	// Initialisation de Fiber
	app := fiber.New()

	// Définir le handler pour le test
	app.Post("/upload", UploadFile)

	// Configuration de la variable d'environnement
	os.Setenv("MINIO_BUCKET", "test-bucket")

	// Mock MinIO client
	mockMinioClient := new(MockMinioClient)
	config.MinioClient = mockMinioClient

	// Préparer le mock pour PutObject avec une erreur
	mockMinioClient.On("PutObject", mock.Anything, "test-bucket", "test-file.txt", mock.Anything, int64(14), minio.PutObjectOptions{ContentType: "text/plain"}).
		Return(minio.UploadInfo{}, assert.AnError)

	// Créer un fichier de test pour l'upload
	fileContent := "test file content"
	fileBuffer := bytes.NewBufferString(fileContent)

	// Créer une requête POST pour uploader le fichier
	req, err := http.NewRequest("POST", "/upload", fileBuffer)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=---boundary")

	// Correction : Utiliser uniquement le nom du fichier pour `FormFile`
	req.MultipartForm = &multipart.Form{
		File: map[string][]*multipart.FileHeader{
			"fileUpload": {
				&multipart.FileHeader{
					Filename: "test-file.txt",
					Size:     int64(len(fileContent)),
				},
			},
		},
	}

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	// Lire la réponse avec ioutil.ReadAll au lieu de utils.Body
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), `"error":true`)
	assert.Contains(t, string(body), `"msg":"upload error"`)
}

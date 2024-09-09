package test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/minio/minio-go"
	"github.com/stretchr/testify/assert"
)

func TestUploadFileHandler(t *testing.T) {
	client, err := initMinioClient()
	if err != nil {
		t.Fatalf("Échec de l'initialisation du client MinIO: %v", err)
	}
	minioClient = client

	// Créez un buffer pour la requête multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "testfile.txt")
	if err != nil {
		t.Fatalf("Échec de la création du formulaire: %v", err)
	}
	part.Write([]byte("test data"))
	writer.WriteField("bucket", "test-bucket")
	writer.Close()

	req, err := http.NewRequest("POST", "/upload-file", body)
	if err != nil {
		t.Fatalf("Échec de la création de la requête: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(uploadFileHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Code de statut attendu %v, mais obtenu %v", http.StatusOK, rr.Code)

	expected := `{"message":"Fichier uploadé avec succès"}`
	assert.JSONEq(t, expected, rr.Body.String(), "Corps de la réponse attendu %v, mais obtenu %v", expected, rr.Body.String())

	// Vérifiez que le fichier a bien été uploadé
	ctx := context.Background()
	object, err := minioClient.GetObject(ctx, "test-bucket", "testfile.txt", minio.GetObjectOptions{})
	if err != nil {
		t.Fatalf("Échec de la récupération du fichier: %v", err)
	}
	object.Close()
}

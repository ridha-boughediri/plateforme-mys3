package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadFileHandler(t *testing.T) {
	client, err := initMinioClient()
	if err != nil {
		t.Fatalf("Échec de l'initialisation du client MinIO: %v", err)
	}
	minioClient = client

	req, err := http.NewRequest("GET", "/download-file?bucket=test-bucket&file=testfile.txt", nil)
	if err != nil {
		t.Fatalf("Échec de la création de la requête: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(downloadFileHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Code de statut attendu %v, mais obtenu %v", http.StatusOK, rr.Code)
}

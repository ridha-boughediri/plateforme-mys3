package test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/minio/minio-go"
	// Remplacez `yourmodulepath` par le chemin de votre module
)

// Initialiser le client MinIO pour les tests
func initMinioClient() (*minio.Client, error) {
	// Votre code pour initialiser le client MinIO
}

// Créer une requête HTTP pour un handler spécifique
func createRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Créer un enregistreur de réponse pour capturer les réponses
func createRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"plateforme-mys3/config"
	"plateforme-mys3/handlers"

	"github.com/gofiber/fiber/v2"
)

func setupApp() *fiber.App {
	app := fiber.New()

	// Initialiser le client MinIO
	err := config.InitMinioClient()
	if err != nil {
		panic("Erreur lors de l'initialisation du client MinIO pour les tests: " + err.Error())
	}

	// Définir les routes pour les handlers
	app.Post("/upload-file", handlers.UploadFile)
	app.Get("/create-bucket", handlers.CreateBucketHandler)
	app.Get("/delete-file", handlers.DeleteFileHandler)
	app.Get("/download-file", handlers.DownloadFileHandler)

	return app
}

func TestUploadFile(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("POST", "/upload-file", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Erreur lors de la requête : %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Statut incorrect : attendu %v, obtenu %v", http.StatusOK, resp.StatusCode)
	}
}

func TestCreateBucket(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("GET", "/create-bucket?bucket=testbucket", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Erreur lors de la requête : %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Statut incorrect : attendu %v, obtenu %v", http.StatusCreated, resp.StatusCode)
	}
}

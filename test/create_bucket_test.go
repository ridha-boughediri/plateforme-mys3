package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"plateforme-mys3/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateBucketHandler(t *testing.T) {
	// Créez une instance de Fiber
	app := fiber.New()

	// Enregistrez le handler
	app.Get("/create-bucket", controllers.CreateBucketHandler)

	// Créez une requête HTTP pour le handler createBucketHandler
	req, err := http.NewRequest("GET", "/create-bucket?bucket=test-bucket", nil)
	if err != nil {
		t.Fatalf("Échec de la création de la requête: %v", err)
	}

	// Créez un enregistreur de réponse pour capturer la réponse du handler
	rr := httptest.NewRecorder()

	// Utilisez Fiber pour envoyer la requête
	app.Test(req)

	// Vérifiez le code de statut
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Vérifiez la réponse
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Erreur de décodage de la réponse: %v", err)
	}
	assert.Equal(t, "Bucket créé avec succès", response["message"])
}

// internal/middleware/auth.go
package middleware

import (
	"log"
	"net/http"
	"plateforme-mys3/config"
	"plateforme-mys3/internal/auth"
)

// AuthMiddleware applique l'authentification AWS SigV4 à tous les handlers
func AuthMiddleware(next http.Handler, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Requête reçue: %s %s", r.Method, r.URL.Path)
		if !auth.VerifyAWSSignature(r, cfg) {
			log.Println("Vérification de la signature échouée")
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`<Error><Code>AccessDenied</Code><Message>Access Denied</Message></Error>`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

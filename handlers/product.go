package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"plateforme-mys3/models" // <-- Utiliser le nom du module correct
	"strconv"

	"github.com/gorilla/mux"
)

// GetProducts récupère tous les produits de la base de données et les renvoie au format JSON.
func GetProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Exécution de la requête SQL pour obtenir tous les produits.
		rows, err := db.Query("SELECT * FROM products")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Initialisation d'une slice pour stocker les produits.
		var products []models.Product

		// Boucle pour parcourir toutes les lignes retournées par la requête.
		for rows.Next() {
			var product models.Product
			// Récupération des données de chaque produit et stockage dans la slice.
			if err := rows.Scan(&product.ID, &product.Name, &product.Quantity); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, product)
		}

		// Encodage des produits en JSON et écriture dans la réponse HTTP.
		json.NewEncoder(w).Encode(products)
	}
}

// CreateProduct ajoute un nouveau produit dans la base de données.
func CreateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product

		// Décodage du corps de la requête JSON pour obtenir les données du produit.
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Insertion du nouveau produit dans la base de données.
		result, err := db.Exec("INSERT INTO products (name, quantity) VALUES (?, ?)", product.Name, product.Quantity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Récupération de l'ID du nouveau produit inséré.
		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Mise à jour de l'ID du produit avec l'ID généré.
		product.ID = int(id)

		// Réponse avec le produit nouvellement créé, au format JSON.
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}

// UpdateProduct met à jour un produit existant dans la base de données par son ID.
func UpdateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product

		// Décodage du corps de la requête JSON pour obtenir les nouvelles données du produit.
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Récupération de l'ID du produit à partir des paramètres de l'URL.
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])

		// Mise à jour des données du produit dans la base de données.
		_, err := db.Exec("UPDATE products SET name = ?, quantity = ? WHERE id = ?", product.Name, product.Quantity, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Réponse HTTP avec le code 204 No Content (pas de contenu).
		w.WriteHeader(http.StatusNoContent)
	}
}

// DeleteProduct supprime un produit de la base de données par son ID.
func DeleteProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupération de l'ID du produit à partir des paramètres de l'URL.
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])

		// Exécution de la requête SQL pour supprimer le produit de la base de données.
		_, err := db.Exec("DELETE FROM products WHERE id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Réponse HTTP avec le code 204 No Content (pas de contenu).
		w.WriteHeader(http.StatusNoContent)
	}
}

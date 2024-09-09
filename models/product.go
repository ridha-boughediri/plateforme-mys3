package models

// Product représente un produit stocké dans l'entrepôt.
type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

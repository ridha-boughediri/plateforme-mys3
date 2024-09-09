package handlers

import (
	"database/sql"
	"net/http"

	"plateforme-mys3/storage"

	"github.com/gin-gonic/gin"
)

// Récupérer tous les produits
func GetProductsGin(c *gin.Context, db *sql.DB) {
	// Logique pour récupérer les produits
	products, err := storage.GetProductsFromDB(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

// Créer un nouveau produit
func CreateProductGin(c *gin.Context, db *sql.DB) {
	var product storage.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := storage.CreateProductInDB(db, product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product": product})
}

// Mettre à jour un produit
func UpdateProductGin(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var product storage.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := storage.UpdateProductInDB(db, id, product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product": product})
}

// Supprimer un produit
func DeleteProductGin(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	if err := storage.DeleteProductFromDB(db, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully", "id": id})
}

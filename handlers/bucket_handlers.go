package handlers

import (
	"database/sql"
	"net/http"

	"plateforme-mys3/storage"

	"github.com/gin-gonic/gin"
)

// Créer un nouveau bucket
func CreateBucketGin(c *gin.Context, db *sql.DB) {
	bucketName := c.Param("name")
	if err := storage.CreateBucketInDB(db, bucketName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bucket created successfully", "bucket": bucketName})
}

// Lister tous les buckets
func ListBucketsGin(c *gin.Context, db *sql.DB) {
	buckets, err := storage.ListBucketsFromDB(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"buckets": buckets})
}

// Supprimer un bucket
func DeleteBucketGin(c *gin.Context, db *sql.DB) {
	bucketName := c.Param("name")
	if err := storage.DeleteBucketFromDB(db, bucketName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bucket deleted successfully", "bucket": bucketName})
}

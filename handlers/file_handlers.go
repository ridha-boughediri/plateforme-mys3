package handlers

import (
	"database/sql"
	"net/http"

	"plateforme-mys3/storage"

	"github.com/gin-gonic/gin"
)

// Uploader un fichier dans un bucket
func UploadFileGin(c *gin.Context, db *sql.DB) {
	bucketName := c.Param("name")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Logique pour uploader le fichier
	err = storage.UploadFileToBucket(db, bucketName, file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file": file.Filename, "bucket": bucketName})
}

// Lister les fichiers dans un bucket
func ListFilesGin(c *gin.Context, db *sql.DB) {
	bucketName := c.Param("name")

	files, err := storage.ListFilesInBucket(db, bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bucket": bucketName, "files": files})
}

// Télécharger un fichier spécifique
func DownloadFileGin(c *gin.Context, db *sql.DB) {
	bucketName := c.Param("name")
	fileName := c.Param("file")

	fileData, err := storage.DownloadFileFromBucket(db, bucketName, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.Writer.Write(fileData)
}

// Supprimer un fichier dans un bucket
func DeleteFileGin(c *gin.Context, db *sql.DB) {
	bucketName := c.Param("name")
	fileName := c.Param("file")

	err := storage.DeleteFileFromBucket(db, bucketName, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully", "file": fileName, "bucket": bucketName})
}

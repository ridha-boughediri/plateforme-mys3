package main

import (
	"context"
	"fmt"
	"log"
	"plateforme-mys3/handlers"
	"plateforme-mys3/storage"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // Importer le pilote MySQL
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// Configuration du DSN pour MySQL
	dsn := "root:@tcp(127.0.0.1:3306)/mys3db?charset=utf8mb4&parseTime=True&loc=Local"

	// Initialiser la base de données MySQL
	db, err := storage.InitDB(dsn)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Migrer les tables nécessaires
	if err := storage.Migrate(db); err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}

	// Configuration du client MinIO
	endpoint := "localhost:9000"
	accessKeyID := "admin"
	secretAccessKey := "admin"
	useSSL := false

	// Initialiser le client MinIO
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Successfully connected to MinIO")

	// Nom du bucket
	bucketName := "mybucket"
	location := "us-east-1"

	// Créer un bucket
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created bucket %s\n", bucketName)
	}

	// Opérations de fichiers (exemple)
	uploadFile(minioClient, bucketName)
	downloadFile(minioClient, bucketName)
	deleteFile(minioClient, bucketName)

	// Créer un routeur Gin
	r := gin.Default()

	// Définir les routes pour les produits (exemple existant)
	r.GET("/products", func(c *gin.Context) {
		handlers.GetProductsGin(c, db)
	})
	r.POST("/products", func(c *gin.Context) {
		handlers.CreateProductGin(c, db)
	})
	r.PUT("/products/:id", func(c *gin.Context) {
		handlers.UpdateProductGin(c, db)
	})
	r.DELETE("/products/:id", func(c *gin.Context) {
		handlers.DeleteProductGin(c, db)
	})

	// Définir les routes pour la gestion des buckets
	r.POST("/buckets/:name", func(c *gin.Context) {
		handlers.CreateBucketGin(c, db)
	})
	r.GET("/buckets", func(c *gin.Context) {
		handlers.ListBucketsGin(c, db)
	})
	r.DELETE("/buckets/:name", func(c *gin.Context) {
		handlers.DeleteBucketGin(c, db)
	})

	// Définir les routes pour la gestion des fichiers dans les buckets
	r.POST("/buckets/:name/upload", func(c *gin.Context) {
		handlers.UploadFileGin(c, db)
	})
	r.GET("/buckets/:name/files", func(c *gin.Context) {
		handlers.ListFilesGin(c, db)
	})
	r.GET("/buckets/:name/files/:file", func(c *gin.Context) {
		handlers.DownloadFileGin(c, db)
	})
	r.DELETE("/buckets/:name/files/:file", func(c *gin.Context) {
		handlers.DeleteFileGin(c, db)
	})

	// Démarrer le serveur
	log.Println("API is running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Unable to start server:", err)
	}
}

// Fonction pour uploader un fichier
func uploadFile(minioClient *minio.Client, bucketName string) {
	objectName := "example.txt"
	filePath := "C:\\Data\\example.txt" // Remplacez par le chemin réel de votre fichier
	contentType := "text/plain"

	info, err := minioClient.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}

// Fonction pour télécharger un fichier
func downloadFile(minioClient *minio.Client, bucketName string) {
	objectName := "example.txt"
	err := minioClient.FGetObject(context.Background(), bucketName, objectName, "C:\\Data\\downloaded_example.txt", minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully downloaded %s\n", objectName)
}

// Fonction pour supprimer un fichier
func deleteFile(minioClient *minio.Client, bucketName string) {
	objectName := "example.txt"
	err := minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully deleted %s\n", objectName)
}

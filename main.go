package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

// Initialiser le client MinIO
func initMinioClient() (*minio.Client, error) {
	// Définir les paramètres de connexion MinIO
	endpoint := "localhost:9000"
	accessKeyID := "admin"
	secretAccessKey := "admin1234@@@test"
	useSSL := false

	// Remplacer les valeurs par défaut par celles des variables d'environnement si elles existent
	if envEndpoint := os.Getenv("AZURE_S3_ENDPOINT"); envEndpoint != "" {
		endpoint = envEndpoint
	}

	if envAccessKey := os.Getenv("MINIO_ROOT_USER"); envAccessKey != "" {
		accessKeyID = envAccessKey
	}

	if envSecretKey := os.Getenv("MINIO_ROOT_PASSWORD"); envSecretKey != "" {
		secretAccessKey = envSecretKey
	}

	if envUseSSL := os.Getenv("AZURE_USE_SSL"); envUseSSL == "true" {
		useSSL = true
	}

	// Initialiser le client MinIO avec les options
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Handler pour créer un bucket
func createBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Query().Get("bucket")
	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bucket created successfully"})
	log.Printf("Bucket %s created successfully", bucketName)
}

// Handler pour lister les buckets
func listBucketsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	buckets, err := minioClient.ListBuckets(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var bucketNames []string
	for _, bucket := range buckets {
		bucketNames = append(bucketNames, bucket.Name)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bucketNames)
	log.Println("Buckets listed successfully")
}

// Handler pour uploader un fichier
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Limite la taille du fichier à 10MB
	r.ParseMultipartForm(10 << 20) // 10MB

	bucketName := r.FormValue("bucket")
	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file from form: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	ctx := r.Context()

	// Vérifier si le bucket existe
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		http.Error(w, "Failed to check bucket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Bucket does not exist", http.StatusNotFound)
		return
	}

	// Upload du fichier vers le bucket
	_, err = minioClient.PutObject(ctx, bucketName, handler.Filename, file, handler.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		http.Error(w, "Failed to upload file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded successfully"})
	log.Printf("File %s uploaded successfully to bucket %s", handler.Filename, bucketName)
}

func main() {
	// Initialiser le client MinIO
	var err error
	minioClient, err = initMinioClient()
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	// Configurer les routes HTTP
	http.HandleFunc("/create-bucket", createBucketHandler)
	http.HandleFunc("/list-buckets", listBucketsHandler)
	http.HandleFunc("/upload-file", uploadFileHandler)

	// Démarrer le serveur HTTP
	log.Println("Server is running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

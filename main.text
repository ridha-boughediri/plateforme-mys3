package main

import (
	"encoding/json"
	"io"
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

	// Remplacer les valeurs par celles des variables d'environnement si elles existent
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
		http.Error(w, "Le nom du bucket est requis", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bucket créé avec succès"})
	log.Printf("Bucket %s créé avec succès", bucketName)
}

// Handler pour uploader un fichier
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Limite la taille du fichier à 10MB
	r.ParseMultipartForm(10 << 20) // 10MB

	bucketName := r.FormValue("bucket")
	if bucketName == "" {
		http.Error(w, "Le nom du bucket est requis", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Impossible de récupérer le fichier: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	ctx := r.Context()

	// Vérifier si le bucket existe
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		http.Error(w, "Impossible de vérifier le bucket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Le bucket n'existe pas", http.StatusNotFound)
		return
	}

	// Upload du fichier vers le bucket
	contentType := handler.Header.Get("Content-Type")
	info, err := minioClient.PutObject(ctx, bucketName, handler.Filename, file, handler.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		http.Error(w, "Échec de l'upload du fichier: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Fichier uploadé avec succès",
		"objectName": info.Key,
		"bucketName": bucketName,
		"size":       info.Size,
	})
	log.Printf("Fichier %s uploadé avec succès dans le bucket %s", handler.Filename, bucketName)
}

// Handler pour télécharger un fichier depuis un bucket
func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Query().Get("bucket")
	objectName := r.URL.Query().Get("file")

	if bucketName == "" || objectName == "" {
		http.Error(w, "Le nom du bucket et du fichier sont requis", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	object, err := minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, "Échec du téléchargement du fichier: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer object.Close()

	// Envoyer le fichier en tant que pièce jointe
	w.Header().Set("Content-Disposition", "attachment; filename="+objectName)
	w.Header().Set("Content-Type", "application/octet-stream")

	if _, err := io.Copy(w, object); err != nil {
		http.Error(w, "Échec de l'envoi du fichier: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Fichier %s téléchargé avec succès depuis le bucket %s", objectName, bucketName)
}

// Handler pour supprimer un fichier d'un bucket
func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Query().Get("bucket")
	objectName := r.URL.Query().Get("file")

	if bucketName == "" || objectName == "" {
		http.Error(w, "Le nom du bucket et du fichier sont requis", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		http.Error(w, "Échec de la suppression du fichier: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Fichier %s supprimé avec succès du bucket %s", objectName, bucketName)
	json.NewEncoder(w).Encode(map[string]string{"message": "Fichier supprimé avec succès"})
}

// Handler pour lister les fichiers dans un bucket
func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Query().Get("bucket")
	if bucketName == "" {
		http.Error(w, "Le nom du bucket est requis", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Vérifier si le bucket existe
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		http.Error(w, "Impossible de vérifier le bucket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Le bucket n'existe pas", http.StatusNotFound)
		return
	}

	// Lister les fichiers
	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})

	var files []map[string]interface{}
	for object := range objectCh {
		if object.Err != nil {
			http.Error(w, "Échec de la liste des fichiers: "+object.Err.Error(), http.StatusInternalServerError)
			return
		}
		files = append(files, map[string]interface{}{
			"key":          object.Key,
			"size":         object.Size,
			"lastModified": object.LastModified,
			"etag":         object.ETag,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
	log.Printf("Fichiers listés avec succès pour le bucket %s", bucketName)
}

func main() {
	// Initialiser le client MinIO
	var err error
	minioClient, err = initMinioClient()
	if err != nil {
		log.Fatalf("Échec de l'initialisation du client MinIO: %v", err)
	}

	// Configurer les routes HTTP
	http.HandleFunc("/create-bucket", createBucketHandler)
	http.HandleFunc("/upload-file", uploadFileHandler)
	http.HandleFunc("/list-files", listFilesHandler)
	http.HandleFunc("/download-file", downloadFileHandler) // Nouveau handler pour télécharger des fichiers
	http.HandleFunc("/delete-file", deleteFileHandler)     // Nouveau handler pour supprimer des fichiers

	// Démarrer le serveur HTTP
	log.Println("Serveur en cours d'exécution sur le port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

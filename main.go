package main

import (
	"context"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func init() {
	var err error
	endpoint := "localhost:9000"    // Remplacez par l'URL de votre serveur MinIO
	accessKeyID := "admin"     // Votre clé d'accès
	secretAccessKey := "admin1234@@@test" // Votre clé secrète
	useSSL := false                 // Mettez à true si vous utilisez HTTPS

	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func createBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Query().Get("bucket")
	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Bucket created successfully"))
}

func listBucketsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	buckets, err := minioClient.ListBuckets(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, bucket := range buckets {
		w.Write([]byte(bucket.Name + "\n"))
	}
}

func main() {
	http.HandleFunc("/create-bucket", createBucketHandler)
	http.HandleFunc("/list-buckets", listBucketsHandler)

	log.Println("Server is running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

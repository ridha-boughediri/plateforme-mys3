package storage

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"example.com/hello/app/dto"
	"github.com/gorilla/mux"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to the storage service\n")
}

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]

	bucketPath := filepath.Join(os.Getenv("BUCKET_PATH"), bucketName)

	err := os.Mkdir(bucketPath, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", r.URL.String())
	w.WriteHeader(http.StatusOK)
}

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	maxBuckets := queryParams.Get("max-buckets")

	if maxBuckets == "" {
		maxBuckets = "1000"
	}

	maxBucketsInt, err := strconv.Atoi(maxBuckets)
	if err != nil {
		http.Error(w, "Invalid max-buckets value", http.StatusBadRequest)
		return
	}

	buckets, err := filepath.Glob(os.Getenv("BUCKET_PATH") + "/*")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ListAllMyBucketsResponse{
		XMLName: xml.Name{
			Local: "ListAllMyBucketsResult",
		},
		Xmlns:   "http://s3.amazonaws.com/doc/2006-03-01/",
		Buckets: make([]dto.ListBuckets, 0),
	}

	for i, bucket := range buckets {
		if i >= maxBucketsInt {
			break
		}

		fileInfo, err := os.Stat(bucket)
		if err != nil {
			fmt.Printf("Error retrieving info for bucket %s: %v", bucket, err)
			continue
		}

		creationTime, _ := time.Parse(time.RFC3339, fileInfo.ModTime().Format(time.RFC3339))

		b := dto.ListBuckets{
			Name:         filepath.Base(bucket),
			CreationDate: creationTime,
		}
		resp.Buckets = append(resp.Buckets, b)
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(resp)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]

	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	bucketPath := filepath.Join(os.Getenv("BUCKET_PATH"), bucketName)

	err := os.RemoveAll(bucketPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

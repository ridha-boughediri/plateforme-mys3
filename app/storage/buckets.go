package storage

import (
	"database/sql"
	"encoding/xml"
	"net/http"
	"time"

	"example.com/hello/app/database"
	"example.com/hello/app/dto"
	"github.com/gorilla/mux"
)

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	rows, err := db.Query("SELECT name, created_at FROM buckets")
	if err != nil {
		http.Error(w, "Failed to list buckets: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	resp := dto.ListAllMyBucketsResponse{
		Xmlns:   "http://s3.amazonaws.com/doc/2006-03-01/",
		Buckets: make([]dto.ListBuckets, 0),
	}

	for rows.Next() {
		var bucket dto.ListBuckets
		err := rows.Scan(&bucket.Name, &bucket.CreationDate)
		if err != nil {
			http.Error(w, "Failed to scan bucket: "+err.Error(), http.StatusInternalServerError)
			return
		}
		resp.Buckets = append(resp.Buckets, bucket)
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(resp)
}

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]

	_, err := db.Exec("INSERT INTO buckets (name) VALUES ($1)", bucketName)
	if err != nil {
		http.Error(w, "Failed to create bucket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]

	// Check if bucket is empty before deleting
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM objects WHERE bucket_id = (SELECT id FROM buckets WHERE name = $1)", bucketName).Scan(&count)
	if err != nil || count > 0 {
		http.Error(w, "Bucket is not empty or error occurred", http.StatusConflict)
		return
	}

	_, err = db.Exec("DELETE FROM buckets WHERE name = $1", bucketName)
	if err != nil {
		http.Error(w, "Failed to delete bucket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

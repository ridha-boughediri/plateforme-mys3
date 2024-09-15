package controller

import (
	"encoding/xml"
	"net/http"
	"example.com/hello/app/dto"
	"example.com/hello/app/service"
	"github.com/gorilla/mux"
)

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	rows, err := service.ListBuckets()
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
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]

	err := service.CreateBucket(bucketName)
	if err != nil {
		http.Error(w, "Failed to create bucket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]

	err := service.DeleteBucket(bucketName)
	if err != nil {
		http.Error(w, "Failed to delete bucket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

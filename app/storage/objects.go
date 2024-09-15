package storage

import (
	"database/sql"
	"encoding/xml"
	"io"
	"net/http"
	"strconv"
	"time"

	"example.com/hello/app/database"
	"example.com/hello/app/dto"
	"github.com/gorilla/mux"
)

func AddObject(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]
	objectName := vars["objectName"]

	bucketID := 0
	err := db.QueryRow("SELECT id FROM buckets WHERE name = $1", bucketName).Scan(&bucketID)
	if err != nil {
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read object data", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO objects (bucket_id, name, data) VALUES ($1, $2, $3)", bucketID, objectName, data)
	if err != nil {
		http.Error(w, "Failed to store object: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Upload successful"))
}

func ListObjects(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]

	bucketID := 0
	err := db.QueryRow("SELECT id FROM buckets WHERE name = $1", bucketName).Scan(&bucketID)
	if err != nil {
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return
	}

	rows, err := db.Query("SELECT name, created_at FROM objects WHERE bucket_id = $1", bucketID)
	if err != nil {
		http.Error(w, "Failed to list objects: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	resp := dto.ListObjectsResponse{
		Xmlns:    "http://s3.amazonaws.com/doc/2006-03-01/",
		Name:     bucketName,
		Contents: make([]dto.Object, 0),
	}

	for rows.Next() {
		var object dto.Object
		err := rows.Scan(&object.Key, &object.LastModified)
		if err != nil {
			http.Error(w, "Failed to scan object: "+err.Error(), http.StatusInternalServerError)
			return
		}
		resp.Contents = append(resp.Contents, object)
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(resp)
}

func DownloadObject(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]
	objectName := vars["objectName"]

	bucketID := 0
	err := db.QueryRow("SELECT id FROM buckets WHERE name = $1", bucketName).Scan(&bucketID)
	if err != nil {
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return
	}

	row := db.QueryRow("SELECT data FROM objects WHERE bucket_id = $1 AND name = $2", bucketID, objectName)
	var data []byte
	err = row.Scan(&data)
	if err != nil {
		http.Error(w, "Object not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
}

func DeleteObject(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]
	objectName := vars["objectName"]

	bucketID := 0
	err := db.QueryRow("SELECT id FROM buckets WHERE name = $1", bucketName).Scan(&bucketID)
	if err != nil {
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return
	}

	_, err = db.Exec("DELETE FROM objects WHERE bucket_id = $1 AND name = $2", bucketID, objectName)
	if err != nil {
		http.Error(w, "Failed to delete object: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

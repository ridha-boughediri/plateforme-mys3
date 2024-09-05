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

func ListObjects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]

	queryParams := r.URL.Query()
	prefix := queryParams.Get("prefix")
	marker := queryParams.Get("marker")
	maxKeys := queryParams.Get("max-keys")

	if maxKeys == "" {
		maxKeys = "1000"
	}

	maxKeysInt, err := strconv.Atoi(maxKeys)
	if err != nil {
		http.Error(w, "Invalid max-keys value", http.StatusBadRequest)
		return
	}

	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	bucketPath := filepath.Join(os.Getenv("BUCKET_PATH"), bucketName)

	objects, err := filepath.Glob(filepath.Join(bucketPath, prefix+"*"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ListObjectsResponse{
		Xmlns:       "http://s3.amazonaws.com/doc/2006-03-01/",
		Name:        bucketName,
		Prefix:      prefix,
		Marker:      marker,
		MaxKeys:     maxKeysInt,
		IsTruncated: false,
		Contents:    make([]dto.Object, 0),
	}

	for i, object := range objects {
		if i >= maxKeysInt {
			resp.IsTruncated = true
			break
		}

		fileInfo, err := os.Stat(object)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		creationTime, _ := time.Parse(time.RFC3339, fileInfo.ModTime().Format(time.RFC3339))

		resp.Contents = append(resp.Contents, dto.Object{
			Key:          filepath.Base(object),
			LastModified: creationTime,
			Size:         int(fileInfo.Size()),
		})
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(resp)
}

func CheckObjectExist(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	bucketName := vars["bucketName"]
	objectName := vars["objectName"]

	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
	}

	if objectName == "" {
		http.Error(w, "Object name is required", http.StatusBadRequest)
	}

	objectPath := filepath.Join(os.Getenv("BUCKET_PATH"), bucketName, objectName)

	_, err := os.Stat(objectPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusNotFound)
	}

	fileInfo, err := os.Stat(objectPath)
	if os.IsNotExist(err) {
		http.Error(w, "Object not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastModified := fileInfo.ModTime().Format(http.TimeFormat)

	w.Header().Set("Last-Modified", lastModified)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	w.WriteHeader(http.StatusOK)
}

func DeleteObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]

	object, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var deleteReq dto.DeleteObjectRequest
	err = xml.Unmarshal(object, &deleteReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	objectName := deleteReq.Object.Key

	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	objectPath := filepath.Join(os.Getenv("BUCKET_PATH"), bucketName, objectName)
	err = os.Remove(objectPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.DeleteResult{
		DeletedResult: []dto.Deleted{
			{
				Key: objectName,
			},
		},
	}

	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(resp)
}

package storage

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"example.com/hello/app/dto"
	"github.com/gorilla/mux"
)

func processChunkedStream(reader io.Reader, writer io.Writer) error {
	bufReader := bufio.NewReader(reader)

	for {
		// Read the chunk size line by line
		line, err := bufReader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading chunk size: %v", err)
		}

		// Trim whitespace and split the chunk size from the chunk-signature
		line = strings.TrimSpace(line)
		parts := strings.SplitN(line, ";", 2)
		chunkSizeHex := parts[0]

		// Convert the hex size to an integer
		chunkSize, err := strconv.ParseInt(chunkSizeHex, 16, 64)
		if err != nil {
			return fmt.Errorf("error parsing chunk size: %v", err)
		}

		// If the chunk size is zero, it's the end of the stream
		if chunkSize == 0 {
			break
		}

		// Read the chunk data based on the chunk size
		if _, err := io.CopyN(writer, bufReader, chunkSize); err != nil {
			return fmt.Errorf("error reading chunk data: %v", err)
		}

		// Discard the trailing CRLF after the chunk data
		if _, err := bufReader.Discard(2); err != nil {
			return fmt.Errorf("error discarding CRLF: %v", err)
		}

		// Optional: Handle or discard the chunk-signature if needed
		if len(parts) > 1 {
			chunkSignature := parts[1]
			log.Printf("Chunk signature: %s", chunkSignature)
		}
	}

	return nil
}

func AddObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]
	objectName := vars["objectName"]

	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	if objectName == "" {
		http.Error(w, "Object name is required", http.StatusBadRequest)
		return
	}

	objectPath := filepath.Join(os.Getenv("BUCKET_PATH"), bucketName, objectName)

	if _, err := os.Stat(objectPath); err == nil {
		suffix := 1
		objectNameWithoutExt := strings.TrimSuffix(objectName, filepath.Ext(objectName))
		newObjectName := objectNameWithoutExt
		for {
			newObjectName = fmt.Sprintf("%s-%d%s", objectNameWithoutExt, suffix, filepath.Ext(objectName))
			newObjectPath := filepath.Join(os.Getenv("BUCKET_PATH"), bucketName, newObjectName)
			if _, err := os.Stat(newObjectPath); os.IsNotExist(err) {
				objectPath = newObjectPath
				break
			}
			suffix++
		}
	}

	file, err := os.Create(objectPath)
	if err != nil {
		http.Error(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to create file: %v", err)
		return
	}
	defer file.Close()

	if r.Header.Get("X-Amz-Content-Sha256") == "STREAMING-AWS4-HMAC-SHA256-PAYLOAD" {
		err = processChunkedStream(r.Body, file)
		if err != nil {
			http.Error(w, "Failed to write chunked data: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to write chunked data: %v", err)
			return
		}
	} else {
		_, err = io.Copy(file, r.Body)
		if err != nil {
			http.Error(w, "Failed to write data: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to write data: %v", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Upload successful"))
	log.Printf("Successfully uploaded file: %s", objectPath)
}

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

func DownloadObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]
	objectName := vars["objectName"]

	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	if objectName == "" {
		http.Error(w, "Object name is required", http.StatusBadRequest)
		return
	}

	objectPath := filepath.Join(os.Getenv("BUCKET_PATH"), bucketName, objectName)

	fileInfo, err := os.Stat(objectPath)
	if err != nil {
		http.Error(w, "File not found: "+err.Error(), http.StatusNotFound)
		return
	}

	file, err := os.Open(objectPath)
	if err != nil {
		http.Error(w, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+objectName+"\"")
	lastModified := fileInfo.ModTime().UTC().Format(time.RFC1123)
	lastModified = strings.Replace(lastModified, "UTC", "GMT", 1)
	w.Header().Set("Last-Modified", lastModified)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to write file content: "+err.Error(), http.StatusInternalServerError)
		return
	}
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

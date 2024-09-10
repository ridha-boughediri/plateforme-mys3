package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"example.com/hello/app/storage"
	"github.com/gorilla/mux"
)

func TestCreateBucketHandler(t *testing.T) {
	os.Setenv("BUCKET_PATH", "./app/buckets")

	router := mux.NewRouter()
	router.HandleFunc("/{bucketName}/", storage.CreateBucket).Methods("PUT")
	req, err := http.NewRequest("PUT", "/testbucket/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestListBucketsHandler(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/", storage.ListBuckets).Methods("GET")

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(storage.ListBuckets)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAddObjectHandler(t *testing.T) {
	os.Setenv("BUCKET_PATH", "./app/buckets")

	router := mux.NewRouter()
	router.HandleFunc("/{bucketName}/{objectName}", storage.AddObject).Methods("POST")

	filePath := "./testobject.txt"
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/testbucket/testobject.txt", bytes.NewReader(buffer))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(filePath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListObjectsHandler(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/{bucketName}/", storage.ListObjects).Methods("GET", "HEAD")

	req, err := http.NewRequest("GET", "/testbucket/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestDownloadObjectHandler(t *testing.T) {
	os.Setenv("BUCKET_PATH", "./app/buckets")

	router := mux.NewRouter()
	router.HandleFunc("/{bucketName}/{objectName}", storage.DownloadObject).Methods("GET")

	req, err := http.NewRequest("GET", "/testbucket/testobject.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	filePath := "./testobject.txt"
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(rr.Body.String())
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteObjectHandler(t *testing.T) {
	os.Setenv("BUCKET_PATH", "./app/buckets")

	router := mux.NewRouter()
	router.HandleFunc("/{bucketName}/", storage.DeleteObject).Methods("POST")

	body := `
			<Delete>
				<Quiet>false</Quiet>
				<Object>
						<Key>testobject.txt</Key>
				</Object>
		</Delete>
	`

	req, err := http.NewRequest("POST", "/testbucket/", bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeleteBucketHandler(t *testing.T) {
	os.Setenv("BUCKET_PATH", "./app/buckets")

	router := mux.NewRouter()
	router.HandleFunc("/{bucketName}/", storage.DeleteBucket).Methods("DELETE")

	req, err := http.NewRequest("DELETE", "/testbucket/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

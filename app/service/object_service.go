package service

import (
	"database/sql"
	"example.com/hello/app/repository"
)

func AddObject(bucketID int, name string, data []byte) error {
	return repository.AddObject(bucketID, name, data)
}

func ListObjects(bucketID int) (*sql.Rows, error) {
	return repository.ListObjects(bucketID)
}

func GetObject(bucketID int, name string) (*sql.Row, error) {
	return repository.GetObject(bucketID, name)
}

func DeleteObject(bucketID int, name string) error {
	return repository.DeleteObject(bucketID, name)
}

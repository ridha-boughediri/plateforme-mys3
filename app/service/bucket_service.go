package service

import (
	"example.com/hello/app/repository"
)

func CreateBucket(name string) error {
	return repository.CreateBucket(name)
}

func ListBuckets() (*sql.Rows, error) {
	return repository.ListBuckets()
}

func DeleteBucket(name string) error {
	isEmpty, err := repository.IsBucketEmpty(name)
	if err != nil || !isEmpty {
		return err
	}
	return repository.DeleteBucket(name)
}

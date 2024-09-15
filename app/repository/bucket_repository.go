package repository

import (
	"database/sql"
	"example.com/hello/app/database"
)

func CreateBucket(name string) error {
	db := database.GetDB()
	_, err := db.Exec("INSERT INTO buckets (name) VALUES ($1)", name)
	return err
}

func ListBuckets() (*sql.Rows, error) {
	db := database.GetDB()
	return db.Query("SELECT name, created_at FROM buckets")
}

func DeleteBucket(name string) error {
	db := database.GetDB()
	_, err := db.Exec("DELETE FROM buckets WHERE name = $1", name)
	return err
}

func IsBucketEmpty(name string) (bool, error) {
	db := database.GetDB()
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM objects WHERE bucket_id = (SELECT id FROM buckets WHERE name = $1)", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

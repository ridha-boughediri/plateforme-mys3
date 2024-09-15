package repository

import (
	"database/sql"
	"example.com/hello/app/database"
)

func AddObject(bucketID int, name string, data []byte) error {
	db := database.GetDB()
	_, err := db.Exec("INSERT INTO objects (bucket_id, name, data) VALUES ($1, $2, $3)", bucketID, name, data)
	return err
}

func ListObjects(bucketID int) (*sql.Rows, error) {
	db := database.GetDB()
	return db.Query("SELECT name, created_at FROM objects WHERE bucket_id = $1", bucketID)
}

func GetObject(bucketID int, name string) (*sql.Row, error) {
	db := database.GetDB()
	return db.QueryRow("SELECT data FROM objects WHERE bucket_id = $1 AND name = $2", bucketID, name), nil
}

func DeleteObject(bucketID int, name string) error {
	db := database.GetDB()
	_, err := db.Exec("DELETE FROM objects WHERE bucket_id = $1 AND name = $2", bucketID, name)
	return err
}

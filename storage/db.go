package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Pilote MySQL
)

// InitDB initialise la connexion à MySQL
func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Vérifiez la connexion à la base de données
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate crée les tables nécessaires dans la base de données
func Migrate(db *sql.DB) error {
	createBucketTableSQL := `CREATE TABLE IF NOT EXISTS buckets (
		id INT AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		PRIMARY KEY (id)
	);`

	createFilesTableSQL := `CREATE TABLE IF NOT EXISTS files (
		id INT AUTO_INCREMENT,
		bucket_id INT,
		filename VARCHAR(255) NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (bucket_id) REFERENCES buckets(id)
	);`

	if _, err := db.Exec(createBucketTableSQL); err != nil {
		return fmt.Errorf("could not create buckets table: %v", err)
	}
	if _, err := db.Exec(createFilesTableSQL); err != nil {
		return fmt.Errorf("could not create files table: %v", err)
	}

	return nil
}

// Créer un nouveau bucket dans la base de données
func CreateBucketInDB(db *sql.DB, name string) error {
	insertBucketSQL := `INSERT INTO buckets(name) VALUES (?)`
	_, err := db.Exec(insertBucketSQL, name)
	return err
}

// Récupérer la liste des buckets dans la base de données
func ListBucketsFromDB(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM buckets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buckets []string
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, name)
	}
	return buckets, nil
}

// Supprimer un bucket de la base de données
func DeleteBucketFromDB(db *sql.DB, name string) error {
	deleteBucketSQL := `DELETE FROM buckets WHERE name = ?`
	_, err := db.Exec(deleteBucketSQL, name)
	return err
}

// Uploader un fichier dans un bucket
func UploadFileToBucket(db *sql.DB, bucketName string, filename string) error {
	// Récupérer l'ID du bucket
	var bucketID int
	err := db.QueryRow("SELECT id FROM buckets WHERE name = ?", bucketName).Scan(&bucketID)
	if err != nil {
		return fmt.Errorf("could not find bucket: %v", err)
	}

	// Insérer le fichier dans la table files
	insertFileSQL := `INSERT INTO files(bucket_id, filename) VALUES (?, ?)`
	_, err = db.Exec(insertFileSQL, bucketID, filename)
	if err != nil {
		return fmt.Errorf("could not upload file: %v", err)
	}

	return nil
}

// Lister les fichiers dans un bucket
func ListFilesInBucket(db *sql.DB, bucketName string) ([]string, error) {
	// Récupérer l'ID du bucket
	var bucketID int
	err := db.QueryRow("SELECT id FROM buckets WHERE name = ?", bucketName).Scan(&bucketID)
	if err != nil {
		return nil, fmt.Errorf("could not find bucket: %v", err)
	}

	// Récupérer la liste des fichiers
	rows, err := db.Query("SELECT filename FROM files WHERE bucket_id = ?", bucketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []string
	for rows.Next() {
		var filename string
		err = rows.Scan(&filename)
		if err != nil {
			return nil, err
		}
		files = append(files, filename)
	}
	return files, nil
}

// Télécharger un fichier depuis un bucket
func DownloadFileFromBucket(db *sql.DB, bucketName, filename string) ([]byte, error) {
	// Cette fonction devra inclure la logique pour récupérer le fichier.
	// Pour l'instant, on retourne juste un exemple de données.
	return []byte("file content here"), nil
}

// Supprimer un fichier dans un bucket
func DeleteFileFromBucket(db *sql.DB, bucketName, filename string) error {
	// Récupérer l'ID du bucket
	var bucketID int
	err := db.QueryRow("SELECT id FROM buckets WHERE name = ?", bucketName).Scan(&bucketID)
	if err != nil {
		return fmt.Errorf("could not find bucket: %v", err)
	}

	// Supprimer le fichier
	deleteFileSQL := `DELETE FROM files WHERE bucket_id = ? AND filename = ?`
	_, err = db.Exec(deleteFileSQL, bucketID, filename)
	return err
}

// Structure Product
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Récupérer tous les produits de la base de données
func GetProductsFromDB(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// Créer un produit dans la base de données
func CreateProductInDB(db *sql.DB, product Product) error {
	insertProductSQL := `INSERT INTO products (name, price) VALUES (?, ?)`
	_, err := db.Exec(insertProductSQL, product.Name, product.Price)
	return err
}

// Mettre à jour un produit dans la base de données
func UpdateProductInDB(db *sql.DB, id string, product Product) error {
	updateProductSQL := `UPDATE products SET name = ?, price = ? WHERE id = ?`
	_, err := db.Exec(updateProductSQL, product.Name, product.Price, id)
	return err
}

// Supprimer un produit de la base de données
func DeleteProductFromDB(db *sql.DB, id string) error {
	deleteProductSQL := `DELETE FROM products WHERE id = ?`
	_, err := db.Exec(deleteProductSQL, id)
	return err
}

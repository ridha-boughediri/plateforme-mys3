package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // Importation du driver SQLite
)

// InitDB initialise la connexion à la base de données SQLite.
func InitDB(filepath string) *sql.DB {
	// Ouverture de la base de données SQLite.
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}

	// Vérification que la connexion à la base de données est bien établie.
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Retourne l'objet DB pour exécuter des requêtes SQL.
	return db
}

// Migrate crée la table "products" si elle n'existe pas déjà.
func Migrate(db *sql.DB) {
	query := `
    CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        quantity INTEGER NOT NULL
    );
    `

	// Exécution de la requête SQL pour créer la table.
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

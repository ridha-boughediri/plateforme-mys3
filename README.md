
# Plateforme-Mys3

## Introduction

Ce projet est une API Server en Go qui utilise MinIO comme solution de stockage compatible avec le protocole S3 pour gérer des buckets et stocker des fichiers. L'API expose des endpoints pour interagir avec MinIO, créant ainsi une plateforme S3 customisée.

## Prérequis

Avant de démarrer, assurez-vous que vous avez installé les éléments suivants :

- **Go** (version 1.18 ou plus récente)
- **MinIO** (pour le serveur de stockage d'objets)
- **MinIO Client (mc)** (pour interagir avec MinIO)
- **Git** (pour cloner le projet)

## Installation

### 1. Cloner le projet

```bash
git clone https://github.com/votre-utilisateur/mon-projet-minio-api.git
cd mon-projet-minio-api
```

### 2. Installer les dépendances

```bash
go mod tidy
```

### 3. Configurer MinIO

#### a. Lancer MinIO

Pour démarrer MinIO, exécutez la commande suivante :

```bash
minio server /chemin/vers/le/répertoire-de-stockage
```

#### b. Définir les identifiants de MinIO

Définissez les variables d'environnement pour les identifiants de MinIO :

```bash
export MINIO_ROOT_USER="admin"
export MINIO_ROOT_PASSWORD="admin1234"
```

### 4. Lancer le serveur Go

```bash
go run main.go
```

Le serveur sera lancé sur le port 3000 par défaut.

## Utilisation

Voici une liste des endpoints disponibles dans l'API :

1. **Créer un Bucket**

   - Endpoint : `/create-bucket`
   - Méthode : `GET`
   - Exemple : 
   
     ```bash
     curl "http://localhost:3000/create-bucket?bucket=mon-bucket"
     ```
   - Description : Ce endpoint permet de créer un bucket dans MinIO.

2. **Lister les Buckets**

   - Endpoint : `/list-buckets`
   - Méthode : `GET`
   - Exemple : 
   
     ```bash
     curl "http://localhost:3000/list-buckets"
     ```
   - Description : Ce endpoint renvoie la liste des buckets disponibles.

3. **Uploader un Fichier**

   - Endpoint : `/upload`
   - Méthode : `POST`
   - Description : Permet de télécharger un fichier dans le bucket spécifié.

## Commandes MinIO Client (`mc`)

Voici quelques commandes de base pour interagir avec votre serveur MinIO via MinIO Client (`mc`).

### a. Configurer l'alias MinIO

```bash
mc alias set myminio http://localhost:9000 admin admin1234
```

### b. Créer un Bucket

```bash
mc mb myminio/mon-bucket
```

### c. Uploader des Fichiers

```bash
mc cp /chemin/vers/fichier.txt myminio/mon-bucket
```

### d. Lister les Fichiers d'un Bucket

```bash
mc ls myminio/mon-bucket
```

### e. Supprimer un Fichier

```bash
mc rm myminio/mon-bucket/fichier.txt
```

---

## Tests

Nous utilisons le package `testify` pour effectuer des tests unitaires sur cette API.

### Installation de `testify`

```bash
go get github.com/stretchr/testify
```

### Exemple de test pour l'API de suppression de fichier

Voici un exemple de test en utilisant `testify` pour tester la fonction de suppression d'un fichier dans un bucket MinIO.

```go
package controllers

import (
	"context"
	"testing"

	"plateforme-mys3/config"

	"github.com/stretchr/testify/assert"
	"github.com/minio/minio-go/v7"
)

func TestDeleteFileHandler(t *testing.T) {
	// Préparer un mock MinIO client
	ctx := context.Background()

	// Remplacez ceci par une configuration mock de MinIO si vous avez
	minioClient := config.MinioClient
	bucketName := "test-bucket"
	fileName := "test-file.txt"

	// Simuler l'ajout d'un fichier dans le bucket
	_, err := minioClient.PutObject(ctx, bucketName, fileName, nil, 0, minio.PutObjectOptions{})
	assert.Nil(t, err)

	// Supprimer le fichier
	err = minioClient.RemoveObject(ctx, bucketName, fileName, minio.RemoveObjectOptions{})
	assert.Nil(t, err)

	// Vérifier si le fichier a été bien supprimé
	_, err = minioClient.StatObject(ctx, bucketName, fileName, minio.StatObjectOptions{})
	assert.NotNil(t, err)
}
```

### Exécution des tests

```bash
go test ./... -v
```

Le bloc de tests ci-dessus permet de valider la suppression correcte des fichiers dans le bucket en utilisant un client MinIO mocké.

---

## Utilisation de Fresh pour le Développement

Pour un développement plus rapide, utilisez l'outil **Fresh** pour surveiller les modifications de fichiers et redémarrer automatiquement l'application.

### Installation de Fresh

```bash
go install github.com/pilu/fresh@latest
```

### Configuration de Fresh

Créez un fichier `runner.conf` dans le répertoire racine de votre projet :

```ini
root: .
cmd: go run main.go
```

### Utilisation de Fresh

Démarrez Fresh avec la commande :

```bash
fresh
```

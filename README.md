# plateforme-mys3
bucket S3 protocole de l’API S3.
API Server GoLand & MinIO Client
Ce projet est une API Server en Go qui utilise MinIO en tant que solution de stockage compatible S3 pour gérer des buckets et stocker des fichiers. L'API expose des endpoints pour créer et interagir avec des buckets dans MinIO.

Prérequis
Avant de démarrer, assurez-vous que vous avez installé les éléments suivants :

Go (version 1.18 ou plus récente)
MinIO pour le serveur de stockage d'objets
MinIO Client (mc) pour interagir avec MinIO
Git pour cloner le projet
Installation
1. Cloner le projet
bash
Copier le code
git clone https://github.com/votre-utilisateur/mon-projet-minio-api.git
cd mon-projet-minio-api
2. Installer les dépendances
Assurez-vous que vous êtes dans le répertoire du projet, puis exécutez la commande suivante pour installer les dépendances :

bash
Copier le code
go mod tidy
3. Configurer MinIO
a. Lancer MinIO
Pour démarrer MinIO, naviguez dans le répertoire où se trouve minio.exe et exécutez la commande suivante :

bash
Copier le code
.\minio.exe server C:\data
Remplacez C:\data par le répertoire où vous souhaitez stocker vos fichiers.

b. Définir les identifiants de MinIO
Assurez-vous de définir vos variables d'environnement pour les identifiants de MinIO avant de démarrer le serveur :

bash
Copier le code
$env:MINIO_ROOT_USER = "admin"
$env:MINIO_ROOT_PASSWORD = "admin1234"
Ou, dans l'invite de commande cmd.exe :

cmd
Copier le code
set MINIO_ROOT_USER=admin
set MINIO_ROOT_PASSWORD=admin1234
4. Lancer le serveur Go
Une fois que MinIO est configuré et en cours d'exécution, vous pouvez démarrer le serveur Go.

5. les commandes client minio


# Commandes MinIO Client (`mc`)

Ce document fournit une liste de commandes MinIO Client (`mc`) pour les opérations courantes sur un serveur MinIO.

## Configurer l'alias du client MinIO

Avant d'effectuer des opérations, vous devez configurer un alias pour votre serveur MinIO. Remplacez `your-access-key` et `your-secret-key` par vos identifiants MinIO réels.

```bash
mc alias set myminio http://localhost:9000 your-access-key your-secret-key
```

## Créer un Bucket

Pour créer un nouveau bucket nommé `larose` :

```bash
mc mb myminio/larose
```

## Uploader des fichiers dans le Bucket

Pour uploader un fichier nommé `your_file.txt` dans le bucket `larose` :

```bash
mc cp /chemin/vers/your_file.txt myminio/larose
```

Pour uploader tous les fichiers d'un répertoire local (`/chemin/vers/repertoire-local`) dans le bucket `larose` :

```bash
mc cp --recursive /chemin/vers/repertoire-local myminio/larose
```

## Lister les fichiers présents dans un Bucket

Pour lister tous les fichiers présents dans le bucket `larose` :

```bash
mc ls myminio/larose
```

Pour une liste plus détaillée, y compris les tailles de fichiers et les dates de modification :

```bash
mc ls --recursive myminio/larose
```

## Télécharger des fichiers spécifiques

Pour télécharger un fichier spécifique nommé `your_file.txt` depuis le bucket `larose` vers le répertoire local actuel :

```bash
mc cp myminio/larose/your_file.txt .
```

Pour télécharger un fichier vers un répertoire local spécifique :

```bash
mc cp myminio/larose/your_file.txt /chemin/vers/repertoire-local/
```

## Supprimer des fichiers spécifiques

Pour supprimer un fichier spécifique nommé `your_file.txt` du bucket `larose` :

```bash
mc rm myminio/larose/your_file.txt
```

Pour supprimer tous les fichiers dans le bucket `larose` (vous serez invité à confirmer) :

```bash
mc rm --recursive --force myminio/larose
```



bash
Copier le code
go run main.go
Le serveur sera lancé sur le port 3000 par défaut.



# Utilisation de Fresh pour le Développement

**Fresh** est un outil utile pour le développement Go qui surveille les modifications dans vos fichiers et redémarre automatiquement votre application lorsque des changements sont détectés. Cela facilite le développement en permettant des itérations rapides sans avoir à redémarrer manuellement le serveur.

## Installation de Fresh

Pour installer `fresh`, utilisez la commande suivante :

```bash
go install github.com/pilu/fresh@latest
```

Assurez-vous que le répertoire `bin` de votre GOPATH est ajouté à votre variable d'environnement `PATH`. Vous pouvez le faire avec la commande suivante (PowerShell) :

```powershell
$env:Path += ";$(go env GOPATH)\bin"
```

## Configuration de Fresh

1. **Créer le Fichier `runner.conf`**

   Dans le répertoire racine de votre projet, créez un fichier nommé `runner.conf` avec le contenu suivant :

   ```ini
   root: .
   cmd: go run main.go
   ```

   - **root:** Définit le répertoire racine de votre projet.
   - **cmd:** La commande que Fresh exécutera pour démarrer votre application.

2. **Utiliser Fresh**

   Avec `fresh` installé et `runner.conf` configuré, vous pouvez démarrer Fresh en exécutant la commande suivante depuis le répertoire racine de votre projet :

   ```bash
   fresh
   ```

   Fresh surveillera automatiquement les fichiers `.go` dans votre projet et redémarrera votre application chaque fois qu'il détecte des modifications.

## Avantages de l'Utilisation de Fresh

- **Développement Rapide**: Modifiez votre code et voyez immédiatement les changements sans redémarrer manuellement votre application.
- **Facile à Configurer**: Une simple configuration avec `runner.conf` suffit pour démarrer.
- **Multiplateforme**: Fonctionne sur Windows, macOS et Linux.



Utilisation
Voici une liste des endpoints disponibles :

1. Créer un Bucket
Endpoint : /create-bucket

Méthode : GET

Exemple de requête :

bash
Copier le code
curl "http://localhost:3000/create-bucket?bucket=mon-bucket"
Description : Ce endpoint permet de créer un bucket dans MinIO.

2. Lister les Buckets
Endpoint : /list-buckets

Méthode : GET

Exemple de requête :

bash
Copier le code
curl "http://localhost:3000/list-buckets"
Description : Ce endpoint renvoie la liste des buckets disponibles dans MinIO.

3. Uploader un Fichier
Endpoint : /upload

Méthode : POST
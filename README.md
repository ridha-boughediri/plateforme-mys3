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

bash
Copier le code
go run main.go
Le serveur sera lancé sur le port 3000 par défaut.

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
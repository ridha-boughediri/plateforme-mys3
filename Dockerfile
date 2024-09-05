# Étape 1 : Utiliser une image Golang pour la construction
FROM golang:1.21 AS builder

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers go.mod et go.sum
COPY go.mod go.sum ./

# Télécharger les dépendances
RUN go mod tidy

# Copier le reste des fichiers de l'application
COPY . .

# Construire l'application
RUN go build -o plateforme-mys3 main.go

# Étape 2 : Créer l'image finale
FROM alpine:latest

# Créer le répertoire de travail
WORKDIR /app

# Copier le binaire depuis l'étape de construction
COPY --from=builder /app/plateforme-mys3 .

# Exposer le port utilisé par l'application
EXPOSE 8080

# Définir la commande à exécuter
CMD ["./plateforme-mys3"]

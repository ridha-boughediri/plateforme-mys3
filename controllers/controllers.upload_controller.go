package controllers

import (
    "context"
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    minioUpload "github.com/minhblues/api/platform/minio"
    "github.com/minio/minio-go/v7"
)

func UploadFile(c *fiber.Ctx) error {
    ctx := context.Background()
    bucketName := os.Getenv("MINIO_BUCKET")
    file, err := c.FormFile("fileUpload")

    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": true,
            "msg":   err.Error(),
        })
    }

    // Get Buffer from file
    buffer, err := file.Open()

    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": true,
            "msg":   err.Error(),
        })
    }
    defer buffer.Close()

     // Create minio connection.
    minioClient, err := minioUpload.MinioConnection()
    if err != nil {
                // Return status 500 and minio connection error.
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": true,
            "msg":   err.Error(),
        })
    }

    objectName := file.Filename
    fileBuffer := buffer
    contentType := file.Header["Content-Type"][0]
    fileSize := file.Size

    // Upload the zip file with PutObject
    info, err := minioClient.PutObject(ctx, bucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": true,
            "msg":   err.Error(),
        })
    }

    log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

    return c.JSON(fiber.Map{
        "error": false,
        "msg":   nil,
        "info":  info,
    })
}

Routes for the API endpoints
// ./pkg/routes/not_found_route.go

package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/minhblues/api/app/controllers"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
    // Create routes group.
    route := a.Group("/api/v1")

    // upload files
    route.Post("/upload", controllers.UploadFile)
}

The main function
package main

import (
    "github.com/gofiber/fiber/v2"
    _ "github.com/joho/godotenv/autoload" // load .env file automatically
    "github.com/minhblues/api/pkg/configs"
    "github.com/minhblues/api/pkg/routes"
    "github.com/minhblues/api/pkg/utils"
)

func main() {
    // Define Fiber config.

    config := configs.FiberConfig()

    // Define a new Fiber app with config.
    app := fiber.New(config)


    routes.PublicRoutes(app)  // Register a public routes for app.


    // Start server (with graceful shutdown).
    utils.StartServer(app)

}


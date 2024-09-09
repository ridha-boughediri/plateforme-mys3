MinhBLues
MinhBLues
Posted on Dec 9, 2021


8

1

2
Easy File Uploads in Go Fiber with Minio
#
go
#
minio
#
beginners
#
tutorial
Introduction
Hello, friends! 😉 Welcome to a really great tutorial. I've tried to make for you as simple step-by-step instructions as possible, based on a real-life application, so that you can apply this knowledge here and now.
I'm writing this tutorial only to share my experience and to show that backend development in Golang using the Fiber framework is easy!

What do we want to build?
Let's create a REST API with fiber which we upload file into Minio.

Setting Minio with Docker
Install and run Docker service for your OS. By the way, in this tutorial I'm using the latest version (at this moment) v20.10.10
docker run \
  -p 9000:9000 \
  -p 9001:9001 \
  --name minio1 \
  -v D:\data:/data \
  -e "MINIO_ROOT_USER=AKIAIOSFODNN7EXAMPLE" \
  -e "MINIO_ROOT_PASSWORD=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" \
  quay.io/minio/minio server /data --console-address ":9001"
☝️ For more information, please visit: https://docs.min.io/docs/minio-docker-quickstart-guide.html

Fiber config in ENV file
# Minio settings:
MINIO_ENDPOINT="localhost:9000"
MINIO_PORT= 9000
MINIO_ACCESSKEY="AKIAIOSFODNN7EXAMPLE"
MINIO_SECRETKEY="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
MINIO_BUCKET="dev-minio"
Minio connection
The minio connection is the most important part of this application.

The method for the connection:
// ./platform/minio/minio.go

package minioUpload

import (
    "context"
    "log"
    "os"

    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioConnection func for opening minio connection.
func MinioConnection() (*minio.Client, error) {
    ctx := context.Background()
    endpoint := os.Getenv("MINIO_ENDPOINT")
    accessKeyID := os.Getenv("MINIO_ACCESSKEY")
    secretAccessKey := os.Getenv("MINIO_SECRETKEY")
    useSSL := false
    // Initialize minio client object.
    minioClient, errInit := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
        Secure: useSSL,
    })
    if errInit != nil {
        log.Fatalln(errInit)
    }

    // Make a new bucket called dev-minio.
    bucketName := os.Getenv("MINIO_BUCKET")
    location := "us-east-1"

    err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
    if err != nil {
        // Check to see if we already own this bucket (which happens if you run this twice)
        exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
        if errBucketExists == nil && exists {
            log.Printf("We already own %s\n", bucketName)
        } else {
            log.Fatalln(err)
        }
    } else {
        log.Printf("Successfully created %s\n", bucketName)
    }
    return minioClient, errInit
}

Create controllers
The principle of the POST methods:

Make a request to the API endpoint;
Parse Form File of request (or an error);
Make a connection to the minio (or an error);
Validate file with a new file from Form-data (or an error);
Upload a new record in the table books (or an error);
Return the status 200 and JSON with a new info file;
// ./app/controllers.upload_controller.go
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

Run project
Some people (including me) crave live reloading in Go, especially the ones who are used to working with interpreted languages like JavaScript, Python, and Ruby. This project I will use nodemon.

☝️ For more ways, please visit: https://techinscribed.com/5-ways-to-live-reloading-go-applications/
nodemon --exec go run main.go --signal SIGNTERM
Test upload file
http://localhost:9001/dashboard:
username: AKIAIOSFODNN7EXAMPLE
password: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
Image description

Postman
http://localhost:5000/api/v1/upload:

Image description

Result Minio
http://localhost:9001/buckets/dev-minio/browse:

Image description

It works. Woohoo! 🎉

P.S.
If you want more articles like this on this blog, then post a comment below and subscribe to me. Thanks! 😘

👋 Before you go

Do your career a favor. Join DEV. (The website you're on right now)

It takes one minute, it's free, and is worth it for your career.

Get started

Top comments (1)
Subscribe
pic
Add to the discussion
 
 
ianmuhia profile image
ian muhia
•
Jan 9 '22

link to the repo please..?


1
 like
Like
Reply
Code of Conduct • Report abuse
Read next
teccmark_ profile image
Why Django Web Development Outshines Traditional CMS Solutions
Teccmark - Aug 18

cudilala profile image
Borrowing and References in Rust Explained
Augustine Madu - Aug 31

pradumnasaraf profile image
Dockerizing a Golang API with MySQL and adding Docker Compose Support
Pradumna Saraf - Sep 4

bgdnvarlamov profile image
The Power of Well-Structured Logs in Software Development
Bogdan Varlamov - Aug 31


MinhBLues
Follow
Joined
Nov 30, 2021
Trending on DEV Community 
Kudzai Murimi profile image
AI Tools That Can Make Your Life as a Web Developer SO Much Easier 😊!
#webdev #devops #discuss #javascript
Rohan Sharma profile image
Top 3 Open-Source Events that will make your October memorable!
#opensource #productivity #learning #github
Kiran Naragund profile image
🤯Powerful AI Tools You Should Know v2🫵
#ai #productivity #programming #tooling
// ./app/controllers.upload_controller.go
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

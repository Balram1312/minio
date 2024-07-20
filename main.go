package main

import (
	"context"
	"log"
	"minio-demo/handler"

	utils "minio-demo/utils"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// corsMiddleware returns a middleware handler that handles CORS requests
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {

	c := gin.Default()

	// Enable CORS middleware
	c.Use(corsMiddleware())

	// MinIO S3 configuration
	ctx := context.Background()

	minioClient := utils.CreateMinioClient(handler.MinioEndpoint, handler.MinioAccessKeyID, handler.MinioSecretAccessKey)
	log.Println("Connection established successfully ...")

	app := handler.MyApp{
		MinioClient: minioClient,
		Context:     ctx,
	}

	app.Routes(c)
	c.Run(":8002")
}

func CreateMinioClient(endpoint string, accessKeyID string, secretAccessKey string) *minio.Client {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln("Error initializing MinIO client:", err)
	}
	return minioClient
}

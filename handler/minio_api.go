package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	urls "net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// The `UploadFile` function is used to upload a file to s3 server.
// it create the key in configured bucket.
//
// - @param file {{file}} : the file parameter used to get the multipart file.
func (app *MyApp) UploadFile(c *gin.Context) {
	// FormFile returns the first file from the request or error if no file is present
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// override bucket name if provided
	InitializeBucketIfProvided(c.Query("bucket"))

	// Read the file
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert file to byte slice
	fileBytes := buffer.Bytes()

	// Example: Print out file information
	log.Printf("File Name: %s\n", header.Filename)
	log.Printf("File Size: %d bytes\n", header.Size)
	log.Printf("MIME Header: %+v\n", header.Header)

	objectName := "test.png"
	reader := bytes.NewReader(fileBytes) // Use bytes1 as the object content
	contentType := "application/octet-stream"
	bucketName := "public-demo-bucket"

	putObject := minio.PutObjectOptions{
		ContentType: contentType,
	}

	info, err := app.MinioClient.PutObject(app.Context, bucketName, header.Filename, reader, int64(len(fileBytes)), putObject)
	if err != nil {
		log.Fatalf("Error uploading object to MinIO: %v", err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	// Return byte array as response (for demonstration purposes)
	c.JSON(http.StatusOK, gin.H{
		"filename": header.Filename,
		"message":  "file uploaded successfully",
	})
}

// The `DownloadFile` method is a function that take gin context as parameter and return server response.
//
// parameters
//
//   - @filename {string}: repsresents the name of file that will be downloaded through its url.
//   - @access_type {string}: repsresents the access type of the file there are 2 conditions,
//     public {string}: repsresents the public access type of the file
//   - @bucket_name {string}: repsresents the bucket name if manually provided.
//
// Returns
//   - Status Code: 200 ||
//     Response
//     {
//     "url":""
//     "error":false
//     }
//   - Status Code: 400 ||
//     Response
//     {
//     "message":""
//     "error":true
//     }
func (app *MyApp) DownloadFile(c *gin.Context) {
	fileName := c.Query("file_name")
	accessType := c.Query("access_type")
	buckerName := c.Query("bucket")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_name query parameter is required"})
		return
	}

	// override bucket name if provided
	InitializeBucketIfProvided(buckerName)

	if strings.EqualFold(accessType, "public") {
		// Construct public object URL
		c.JSON(200, gin.H{
			"url":   CreatePublicUrl(fileName),
			"error": false,
		})
		return
	} else {
		// Gernerate presigned get object url.
		reqParams := make(urls.Values)
		presignedURL, err := CreatePresignedUrl(app, fileName, reqParams)
		if err != nil {
			var err = fmt.Errorf("error generating presigned URL: %s", err.Error())
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"url":   presignedURL.String(),
			"error": false,
		})
	}
}

// ListBuckets retrieves a list of all available buckets in the MinIO server.
// It uses the MinIO client to list all buckets and returns them as a JSON response.
//
// Parameters:
//   - c: A pointer to a gin.Context object, which provides the request context and response writer.
//
// Returns:
//   - Status Code: 200 (http.StatusOK) ||
//     Response:
//     {
//     "buckets": []string
//     }
func (app *MyApp) ListBuckets(c *gin.Context) {
	miniobuckets, err := app.MinioClient.ListBuckets(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	var buckets []string
	for _, bucket := range miniobuckets {
		fmt.Println(bucket)
		buckets = append(buckets, bucket.Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"buckets": buckets,
	})
}

func InitializeBucketIfProvided(buckerName string) {
	if buckerName != "" {
		BucketName = buckerName
	}
}

func CreatePresignedUrl(app *MyApp, fileName string, reqParams urls.Values) (*urls.URL, error) {
	presignedURL, err := app.MinioClient.PresignedGetObject(context.Background(), BucketName, fileName, time.Duration(1000)*time.Second, reqParams)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	log.Printf("presigned URL created: %s", presignedURL.String())
	return presignedURL, nil
}

func CreatePublicUrl(fileName string) string {
	return fmt.Sprintf("%s/%s/%s", MinioEndpoint, BucketName, fileName)
}

// The `DownloadFile` method is a function that take gin context as parameter and return server response.
//
// parameters
//
//   - @filename {string}: repsresents the name of file that will be downloaded through its url.
//   - @access_type {string}: repsresents the access type of the file there are 2 conditions,
//     public {string}: repsresents the public access type of the file
//   - @bucket_name {string}: repsresents the bucket name if manually provided.
//
// Returns
//   - Status Code: 200 ||
//     Response
//     {
//     "url":""
//     "error":false
//     }
//   - Status Code: 400 ||
//     Response
//     {
//     "message":""
//     "error":true
//     }
func (app *MyApp) ListObjects(c *gin.Context) {

	// override bucket name if provided
	InitializeBucketIfProvided(c.Query("bucket"))

	objectCh := app.MinioClient.ListObjects(app.Context, BucketName, minio.ListObjectsOptions{
		// Prefix:    "myprefix",
		Recursive: true,
	})

	var objects []string
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			c.JSON(400, gin.H{
				"error": false,
				"data":  object.Err.Error(),
			})
			return
		}
		log.Println("object: \n", object.Key)
		objects = append(objects, object.Key)
	}
	c.JSON(200, gin.H{
		"error": false,
		"data":  objects,
	})
}

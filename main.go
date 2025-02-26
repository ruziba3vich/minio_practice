package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	bucketName = "uploads"
	minioURL   = "localhost:9000"
	accessKey  = "admin"
	secretKey  = "secretpass"
)

func main() {
	r := gin.Default()
	r.POST("/upload", uploadFile)

	log.Println("Server running on :8080")
	r.Run(":8080")
}

func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	// Open file
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open file"})
		return
	}
	defer f.Close()

	// Generate unique filename
	filename := fmt.Sprintf("%d-%s", time.Now().UnixMilli(), file.Filename)

	// Initialize MinIO client
	minioClient, err := minio.New(minioURL, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MinIO"})
		return
	}

	// Ensure bucket exists
	err = createBucket(minioClient, bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/check bucket"})
		return
	}

	// Upload file to MinIO
	_, err = minioClient.PutObject(
		context.Background(),
		bucketName,
		filename,
		f,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Return file URL
	fileURL := fmt.Sprintf("http://%s/%s/%s", minioURL, bucketName, filename)
	c.JSON(http.StatusOK, gin.H{"url": fileURL})
}

func createBucket(client *minio.Client, bucket string) error {
	exists, err := client.BucketExists(context.Background(), bucket)
	if err != nil {
		return err
	}

	if !exists {
		return client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
	}
	return nil
}

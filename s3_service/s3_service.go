package s3_service

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/sai-subramanian/21BCE0040_Backend.git/configl"
	"github.com/sai-subramanian/21BCE0040_Backend.git/models"
)

type AWSService struct {
	S3Client *s3.Client
}

func AwsInit() (AWSService, error) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		log.Println("error loading default config")
	}
	awsService := AWSService{
		S3Client: s3.NewFromConfig(config),
	}
	return awsService, nil
}

func (awsSvc AWSService) UploadFile(c *gin.Context, req S3Dto) error {
	bucketName := "golang-test-zztzz" // to be moved to config file

	// Open the file
	file, err := os.Open(req.FileName)
	if err != nil {
		log.Println("error while opening the file", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open file",
		})
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("error while closing the file:", err)
		}
	}()

	// Upload the file to the bucket
	_, err = awsSvc.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(req.BucketKey),
		Body:   file,
		ContentType: aws.String(req.ContentType),
	})
	if err != nil {
		log.Println("error while uploading the file to S3", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload file to S3",
		})
		return err
	}

	// Create a new db entry in the table file, with metadata about the file
	newFileDb := models.File{
		Createdby: req.UserId,
		Key:       req.BucketKey,
		ExpirationDate: req.ExpirationDate,
	}
	result := configl.DB.Create(&newFileDb)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create file record in database",
		})
		return result.Error
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded and record created successfully",
	})

	return nil
}

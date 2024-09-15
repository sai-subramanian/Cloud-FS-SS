package s3_service

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

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
	bucketName := os.Getenv("bucketName")

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

	// Uploading the file to the bucket
	_, err = awsSvc.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(req.BucketKey),
		Body:        file,
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
		Createdby:      req.UserId,
		Key:            req.BucketKey,
		ExpirationDate: req.ExpirationDate,
	}
	url, err := awsSvc.GetSignedUrl(req.BucketKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate signed URL",
		})
		return err
	}
	newFileDb.Url = url
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
		"file":    newFileDb,
	})

	return nil
}

func (awsSvc *AWSService) GetSignedUrlHandler(c *gin.Context) {

	objectKey := c.Param("file_id")

	signedUrl, err := awsSvc.GetSignedUrl(objectKey)
	if err != nil {
		log.Println("Failed to generate signed URL:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate signed URL",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"signed_url": signedUrl,
	})
}



func GetFilesByUserId(c *gin.Context) {

	userId := c.Param("userId")

	var files []models.File

	result := configl.DB.Where("createdby = ?", userId).Find(&files)

	if result.Error != nil {
		log.Println("error while fetching files:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve files",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No files found for the given user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}


func (awsSvc *AWSService) GetSignedUrl(objectKey string) (string, error) {
	bucketName := os.Getenv("bucketName")
	// setting the url valid for 5 days in the future when the url expires we can call this endpoint again and
	// reset the db url back with the fresh url
	// or if we want permanent url we can have it using : https://<bucket-name>.s3.amazonaws.com/<object-key>
	// currently i have made the s3 bucket public for time being so that it does not give ACL issue
	// in future we can make it restricted and just return the signed url from here
	presignDuration := time.Duration(5 * 24 * time.Hour)

	presigner := s3.NewPresignClient(awsSvc.S3Client)

	req, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(presignDuration))

	if err != nil {
		log.Println("Error generating presigned URL:", err)
		return "", err
	}

	return req.URL, nil
}

package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sai-subramanian/21BCE0040_Backend.git/s3_service"
)

func FileRoutes(router *gin.Engine,awsService s3_service.AWSService) {
	
	router.GET("/ping",func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})


    // Route to upload a file
	router.POST("/upload", func(c *gin.Context) {
		
		bucketKey := "test2.txt"
		fileName := "/home/sai/Data/projects/21BCE0040_Backend/TextFile(1).txt"

		// Use the s3Client to upload the file to S3
		err := awsService.UploadFile( bucketKey, fileName)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "File upload failed",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "File uploaded successfully",
			})
		}
	})
}
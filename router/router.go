package router

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/sai-subramanian/21BCE0040_Backend.git/s3_service"
	"github.com/sai-subramanian/21BCE0040_Backend.git/user"
)

func FileRoutes(router *gin.Engine,awsService s3_service.AWSService) {
	
	router.GET("/ping",func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})


    // Route to upload a file
	router.POST("/upload", func(c *gin.Context) {
		var req s3_service.S3Dto
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}
		
		currentTime := time.Now()
		tenDaysFromNow := currentTime.AddDate(0, 0, 10)
		req.ExpirationDate = tenDaysFromNow

		if err := awsService.UploadFile(c, req);  err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}
	
		// c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
	})

	router.POST("/signup",user.SignUp )
	router.POST("/login",user.Login )
	
}
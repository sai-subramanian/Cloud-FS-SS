package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sai-subramanian/21BCE0040_Backend.git/configl"
	"github.com/sai-subramanian/21BCE0040_Backend.git/models"
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
	router.GET("/files/:userId",awsService.GetSignedUrlHandler)
	router.GET("/share/:file_id",awsService.GetSignedUrlHandler)

	router.POST("/register",user.SignUp )
	router.POST("/login",user.Login )
	
	router.GET("/search", func(c *gin.Context) {
		
		userId := c.Query("userId")
		fileName := c.Query("fileName")
		
		//  format to be passed : YYYY-MM-DD
		startDate := c.Query("startDate") 
		endDate := c.Query("endDate")     
		fileType := c.Query("fileType")
	
		var files []models.File
		query := configl.DB.Where("createdby = ?", userId)
	
		
		if fileName != "" {
			query = query.Where("key LIKE ?", "%"+fileName+"%")
		}
	
		
		if fileType != "" {
			query = query.Where("key LIKE ?", "%."+fileType)
		}
	
		
		if startDate != "" && endDate != "" {
			query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}
	
		
		result := query.Find(&files)
	
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch files",
			})
			return
		}
	
		
		c.JSON(http.StatusOK, files)
	})
}
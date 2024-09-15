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


   
	router.POST("/upload", func(c *gin.Context) {
   
    file, _, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to retrieve file",
        })
        return
    }
    
    userId := c.PostForm("userId")
    bucketKey := c.PostForm("bucketKey")
    contentType := c.PostForm("contentType")

    currentTime := time.Now()
    tenDaysFromNow := currentTime.AddDate(0, 0, 5)

    
    req := s3_service.S3Dto{
        File:    file, 
        UserId:      userId,
        BucketKey:   bucketKey,
        ContentType: contentType,
        ExpirationDate: tenDaysFromNow,
    }

    
    if err := awsService.UploadFile(c, req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
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
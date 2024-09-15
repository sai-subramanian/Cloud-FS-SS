package s3_service

import (
	"time"

	"github.com/gin-gonic/gin"
)

type S3Dto struct {
	*gin.Context
	
	UserId string
	BucketKey string
	FileName string
	ContentType string
	ExpirationDate time.Time
	
}

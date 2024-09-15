package s3_service

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type S3Dto struct {
	*gin.Context

	UserId         string
	BucketKey      string
	File           io.Reader
	ContentType    string
	ExpirationDate time.Time
}

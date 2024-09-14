package s3_service

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSService struct{
	S3Client *s3.Client
}


func AwsInit() (AWSService,error) {
	config,err :=  config.LoadDefaultConfig(context.TODO(),config.WithRegion("ap-south-1"))
	if(err != nil){
		log.Println("error loading default config")
	}
	awsService := AWSService{
		S3Client: s3.NewFromConfig(config),
	}
	return awsService,nil
}

func (awsSvc AWSService) UploadFile(bucketName string , bucketKey string, fileName string) error{
	file, err := os.Open(fileName) 
	if err != nil{
		log.Println("error while opening the file", err)
	}else{
		defer file.Close()

		_,err  := awsSvc.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
            Key:    aws.String(bucketKey),
            Body:   file,
		})
		if err != nil{
            log.Println("error while uploading the file to S3", err)
		}
	}
	return err
}


package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sai-subramanian/21BCE0040_Backend.git/configl"
	"github.com/sai-subramanian/21BCE0040_Backend.git/router"
	"github.com/sai-subramanian/21BCE0040_Backend.git/s3_service"
)


func init(){
	configl.LoadEnv()
	configl.ConnectDatabase()
}


func main() {
	
	//intilize AWS service for which fxn is written in s3_service
	awsSvc, err := s3_service.AwsInit()
	if err != nil {
		log.Fatalf("Failed to initialize AWS service: %v", err)
	}
	
	

    // Create a new Gin router
    r := gin.Default()
	
    // Register routes
    router.FileRoutes(r,awsSvc)

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

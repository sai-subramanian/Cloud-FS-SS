package main

import (
	"fmt"
	"log"
	
	"github.com/sai-subramanian/21BCE0040_Backend.git/config"
	"github.com/sai-subramanian/21BCE0040_Backend.git/models"
	"gorm.io/gorm"
)

func init() {
	config.LoadEnv()
	config.ConnectDatabase()
}

func main(){
	
	
	
	if config.DB == nil{
		log.Fatal("DB not initialized")
	}

	config.DB.AutoMigrate((&models.User{}),(&models.File{}))

	Session := config.DB.Session(&gorm.Session{PrepareStmt: false})
	
	if Session != nil {
		fmt.Println("Migration Successful")
	}
}
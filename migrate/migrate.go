package main

import (
	"fmt"
	"log"
	
	"github.com/sai-subramanian/21BCE0040_Backend.git/configl"
	"github.com/sai-subramanian/21BCE0040_Backend.git/models"
	"gorm.io/gorm"
)

func init() {
	configl.LoadEnv()
	configl.ConnectDatabase()
}

func main(){
	
	
	
	if configl.DB == nil{
		log.Fatal("DB not initialized")
	}

	configl.DB.AutoMigrate((&models.User{}),(&models.File{}))

	Session := configl.DB.Session(&gorm.Session{PrepareStmt: false})
	
	if Session != nil {
		fmt.Println("Migration Successful")
	}
}
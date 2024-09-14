package config

import (
	// "os"

	"log"
	"os"

	
	"github.com/supabase-community/postgrest-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Client *postgrest.Client

var DB *gorm.DB

func ConnectDatabase(){
	
	 var err error
	dsn := os.Getenv("DSN")

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	  }), &gorm.Config{})

    // DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 	PrepareStmt: false,
	// })

	if err != nil {
		log.Fatal("Error connecting to database");
	} 

	log.Println("Database connected!!!", DB)
}
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Database			string
	ConnectionString	string
	SecretCode			string
)

func init() {
	err := godotenv.Load()

  	if err != nil {
    	log.Fatal("Error loading .env file")
  	}

	Database = os.Getenv("DATABASE")
	ConnectionString = os.Getenv("DATABASE_URL")
	SecretCode = os.Getenv("SECRET_CODE")
}
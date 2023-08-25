package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var (
	Database			string
	ConnectionString	string
	SecretCode			string

	Port				string

	UploadPath			string
)

func init() {
	UploadPath = "webdata"
	Database = "mysql"
	Port = "9003"


	// err := godotenv.Load()

  	// if err != nil {
    // 	log.Fatal("Error loading .env file")
  	// }

	// Database = os.Getenv("DATABASE")
	// ConnectionString = os.Getenv("DATABASE_URL")
	// SecretCode = os.Getenv("SECRET_CODE")


	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if value := viper.Get("connectionString"); value != nil {
		ConnectionString = value.(string)
	}

	if value := viper.Get("Database"); value != nil {
		Database = value.(string)
	}

	if value := viper.Get("SecretCode"); value != nil {
		SecretCode = value.(string)
	}

	if value := viper.Get("port"); value != nil {
		Port = value.(string)
	}

	if value := viper.Get("uploadPath"); value != nil {
		UploadPath = value.(string)
	}


	log.Println(ConnectionString)
}
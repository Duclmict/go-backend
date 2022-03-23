package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// prefix APP

var (
	// APP information
	App_Name string
	App_Environment string	= "Develop"			 //  [Develop, Stagging, Production]
	App_Debug	    string	= "True"			 //  [True, False]
	App_URL		    string   

	// time
	App_TimeZone	string 	= "Utc"				// default [Utc]
	App_Locale		string 	= "ja"				// default [ja]

	// Encryption Key
	App_KEY			string
	App_CIPHER		string	= "AES-256-CBC"		// default	[AES-256-CBC]
)

func LoadENV() {
	errEnv := godotenv.Load()
    if errEnv != nil {
        panic("Failed to load env file")
    }

	// APP
	App_Name 				= os.Getenv("APP_NAME")
	App_Environment 		= os.Getenv("APP_MODE")
	App_Debug 				= os.Getenv("APP_DEBUG")
	App_URL 				= os.Getenv("APP_URL")
	App_KEY					= os.Getenv("APP_KEY")

	// Database
	DB_CONNECT				= os.Getenv("DB_CONNECTION")
	DB_MYSQL_HOST 			= os.Getenv("DB_HOST")
	DB_MYSQL_PORT 			= os.Getenv("DB_PORT")
	DB_MYSQL_DATABASE 		= os.Getenv("DB_DATABASE")
	DB_MYSQL_USER  			= os.Getenv("DB_USERNAME")
	DB_MYSQL_PASSWORD  		= os.Getenv("DB_PASSWORD")
	// DB_MYSQL_UNIX_SOCKET  	= os.Getenv("APP_NAME")

	// Log
	LOG_FOLDER				= os.Getenv("LOG_FOLDER")

	// Mail
	MAIL_DRIVER				= os.Getenv("MAIL_DRIVER")
	MAIL_HOST 				= os.Getenv("MAIL_HOST")
	MAIL_PORT,errEnv  		= strconv.Atoi(os.Getenv("MAIL_PORT"))
	if errEnv != nil {
        panic("Failed to load env file")
    }
	MAIL_USERNAME 			= os.Getenv("MAIL_USERNAME")
	MAIL_PASSWORD 			= os.Getenv("MAIL_PASSWORD")
	MAIL_ENCRYPTION 		= os.Getenv("MAIL_ENCRYPTION")
	MAIL_FROM_ADDRESS 		= os.Getenv("MAIL_FROM_ADDRESS")
	MAIL_FROM_NAME 			= os.Getenv("MAIL_FROM_NAME")

	MAIL_TEMPLATE = map[string]*MailTemplate {
		"hello":  {
			From: 		MAIL_FROM_ADDRESS,
			Name: 		MAIL_FROM_NAME,
			Subject:	"HELLO SEND MAIL",
		},
	}

	// Service AWS
	AWS_ACCESS_KEY_ID		= os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY 	= os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_DEFAULT_REGION 		= os.Getenv("AWS_DEFAULT_REGION")
	AWS_BUCKET 				= os.Getenv("AWS_BUCKET")

}
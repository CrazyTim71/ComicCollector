package env

import (
	"embed"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const VERSION = "v0.1"
const MaxImageFileSize = 4 << 20 // 4 MiB

// default values as fallback
var SERVER_HOST = "127.0.0.1"
var SERVER_PORT = "8080"
var ENV_FILE_LOCATION = ".env"
var MONGODB_DBNAME = "ComicCollector"
var MONGODB_URI = ""
var RSA_FILENAME = "rsa_private_key.pem"
var SIGNUP_ENABLED = false
var TIMEZONE = "Europe/Berlin"

var FrontendFiles embed.FS
var Timezone *time.Location

func InitEnvironment() bool {
	if err := godotenv.Load(ENV_FILE_LOCATION); err != nil {
		log.Fatal("No .env file found")
	}

	return true
}

func InitTimezone() {
	loc, _ := time.LoadLocation(GetTimezone())
	Timezone = loc
}

func GetServerAddress() string {
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")

	if host != "" && port != "" {
		SERVER_HOST = host
		SERVER_PORT = port
	} else {
		log.Println("No address set in .env file")
		log.Println("Defaulting to " + SERVER_HOST + ":" + SERVER_PORT)
	}

	return SERVER_HOST + ":" + SERVER_PORT
}

func GetDatabaseURI() string {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Println("No 'MONGODB_URI' variable set in .env file")
		log.Println("Defaulting to " + MONGODB_URI)
	} else {
		MONGODB_URI = uri
	}

	return MONGODB_URI
}

func GetDatabaseName() string {
	dbName := os.Getenv("MONGODB_DBNAME")
	if dbName == "" {
		log.Println("No 'MONGODB_DBNAME' variable set in .env file")
		log.Println("Defaulting to " + MONGODB_DBNAME)
	} else {
		MONGODB_DBNAME = dbName
	}

	return MONGODB_DBNAME
}

func GetRSAFilename() string {
	rsaFilename := os.Getenv("RSA_FILENAME")
	if rsaFilename == "" {
		log.Println("No 'RSA_FILENAME' variable set in .env file")
		log.Println("Defaulting to " + RSA_FILENAME)
	} else {
		RSA_FILENAME = rsaFilename
	}

	return RSA_FILENAME
}

func GetSignupEnabled() bool {
	signupEnabled := os.Getenv("SIGNUP_ENABLED")
	if signupEnabled == "" {
		log.Println("No 'SIGNUP_ENABLED' variable set in .env file")
		log.Println("Defaulting to false")
		SIGNUP_ENABLED = false
	} else if signupEnabled == "true" {
		SIGNUP_ENABLED = true
	}

	return SIGNUP_ENABLED
}

func GetTimezone() string {
	timezone := os.Getenv("TIMEZONE")
	if timezone == "" {
		log.Println("No 'TIMEZONE' variable set in .env file")
		log.Println("Defaulting to " + TIMEZONE)
	} else {
		TIMEZONE = timezone
	}

	return TIMEZONE
}

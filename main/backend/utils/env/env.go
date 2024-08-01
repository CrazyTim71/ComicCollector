package env

import (
	"embed"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const VERSION = "v0.1"

var SERVER_HOST = "127.0.0.1"
var SERVER_PORT = "8080"
var ENV_FILE_LOCATION = ".env"
var MONGODB_DBNAME = "ComicCollector"
var MONGODB_URI = ""
var RSA_FILENAME = "rsa_private_key.pem"

var Files embed.FS

func InitEnvironment() bool {
	if err := godotenv.Load(ENV_FILE_LOCATION); err != nil {
		log.Fatal("No .env file found")
	}

	return true
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

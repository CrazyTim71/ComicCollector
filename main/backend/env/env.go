package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var SERVER_HOST = "127.0.0.1"
var SERVER_PORT = "8080"
var ENV_FILE_LOCATION = ".env"
var MONGODB_DBNAME = "ComicCollector"
var MONGODB_URI = ""

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
		return ""
	}

	MONGODB_URI = uri
	return MONGODB_URI
}

func GetDatabaseName() string {
	dbName := os.Getenv("MONGODB_DBNAME")
	if dbName == "" {
		log.Println("No 'MONGODB_DBNAME' variable set in .env file")
		log.Println("Defaulting to " + MONGODB_DBNAME)
	}

	MONGODB_DBNAME = dbName
	return MONGODB_DBNAME
}

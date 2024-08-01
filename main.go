package main

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/router"
	"ComicCollector/main/backend/utils/crypt"
	"ComicCollector/main/backend/utils/env"
	"embed"
	"github.com/gin-gonic/gin"
	"log"
)

//go:embed main/frontend/*
var Files embed.FS

func main() {
	// TODO: check the log level & why is the output red in Goland ??
	// TODO: set the log level of gin based on a env flag
	log.Println("\nComicCollector" + "\nVersion: " + env.VERSION)

	// read the embedded files
	_, err := Files.ReadDir("main/frontend")
	if err != nil {
		log.Println("Failed to read the frontend files")
		return
	}
	env.Files = Files

	// init the environment
	env.InitEnvironment()

	// init the db
	database.InitDatabase()

	// load the RSA key
	crypt.InitRSAKey()

	// create the router
	r := gin.Default()
	router.InitRouter(r)

	// start the server
	address := env.GetServerAddress()

	err = r.Run(address)
	if err != nil {
		log.Fatal(err)
		return
	}
}

package main

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/env"
	"ComicCollector/main/backend/router"
	"ComicCollector/main/backend/utils/crypt"
	"embed"
	"github.com/gin-gonic/gin"
	"log"
)

//go:embed main/frontend/*
var Files embed.FS

func main() {
	// read the embedded files
	_, err := Files.ReadDir("main/frontend")
	if err != nil {
		log.Println("Failed to read the frontend files")
		return
	}

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

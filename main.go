package main

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/router"
	"ComicCollector/main/backend/setup"
	"ComicCollector/main/backend/utils/crypt"
	"ComicCollector/main/backend/utils/env"
	"embed"
	"github.com/gin-gonic/gin"
	"log"
)

//go:embed main/frontend/dist*
var FrontendAssets embed.FS

func main() {
	// TODO: check the log level & why is the output red in Goland ??
	// TODO: set the log level of gin based on a env flag
	log.Println("\nComicCollector" + "\nVersion: " + env.VERSION)

	env.FrontendFiles = FrontendAssets

	// init the environment
	env.InitEnvironment()

	// init the timezone
	env.InitTimezone()

	// init the db
	database.InitDatabase()

	// load the RSA key
	crypt.InitRSAKey()

	// check if the database already exists or if this is the first run
	if !database.HasCollection(database.MongoDB, "user") {
		log.Println("Detected the first startup")
		log.Println("Initializing the database and creating the basic users, roles and permissions")

		err := setup.PerformFirstRunTasks()
		if err != nil {
			log.Println("An error occurred while performing the database initialization: ")
			log.Fatalln(err)
		}
	}

	// create the router
	r := gin.Default()
	router.InitBackendRoutes(r)
	router.InitFrontendRoutes(r)

	// start the server
	address := env.GetServerAddress()

	err := r.Run(address)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// TODO: settings to save things like the libraryName, isSignUpEnabled, ...

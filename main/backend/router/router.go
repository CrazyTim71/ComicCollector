package router

import (
	v1 "ComicCollector/main/backend/api/v1"
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/middleware"
	"ComicCollector/main/backend/utils"
	"ComicCollector/main/backend/utils/env"
	"ComicCollector/main/backend/utils/webcontext"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"log"
	"net/http"
)

func InitStaticAssets(r *gin.Engine) bool {
	// add route for the favicon
	r.StaticFileFS("/favicon.ico", "main/frontend/static/favicon.ico", http.FS(env.Files))

	return true
}

func InitFrontendRoutes(r *gin.Engine) bool {
	r.GET("/", func(c *gin.Context) {
		authCookie, err := c.Cookie("auth_token")
		if authCookie != "" && err == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
			return
		}

		SignupEnabled := env.GetSignupEnabled()
		templateSite := template.Must(
			template.ParseFS(
				env.Files,
				"main/frontend/public/index.gohtml",
				"main/frontend/templates/base.gohtml"))

		data := map[string]interface{}{
			"SIGNUP_ENABLED": SignupEnabled,
		}

		err = templateSite.Execute(c.Writer, data)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
			return
		}
	})

	r.GET("/login", func(c *gin.Context) {
		authCookie, err := c.Cookie("auth_token")
		if authCookie != "" && err == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
			return
		}

		SignupEnabled := env.GetSignupEnabled()
		data := map[string]interface{}{
			"SIGNUP_ENABLED": SignupEnabled,
		}

		templateSite := template.Must(
			template.ParseFS(
				env.Files,
				"main/frontend/public/login/index.gohtml",
				"main/frontend/templates/base.gohtml"))

		err = templateSite.Execute(c.Writer, data)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
			return
		}
	})

	r.GET("/register", func(c *gin.Context) {
		authCookie, err := c.Cookie("auth_token")
		if authCookie != "" && err == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
			return
		}

		if !env.GetSignupEnabled() {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		templateSite := template.Must(
			template.ParseFS(
				env.Files,
				"main/frontend/public/register/index.gohtml",
				"main/frontend/templates/base.gohtml"))

		err = templateSite.Execute(c.Writer, nil)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
			return
		}
	})

	r.GET("/dashboard", middleware.JWTAuth(), func(c *gin.Context) {
		// get the userId
		// because of middleware.JWTAuth() we can safely assume that the user id logged in
		userId, err := webcontext.GetUserId(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			return
		}

		// check if the user is an admin
		isAdmin, err := groups.CheckUserGroup(userId, groups.Administrator)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
			return
		}

		// get the username
		user, err := operations.GetOneById[models.User](database.Tables.User, userId)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
			return
		}

		data := map[string]interface{}{
			"isAdmin":  isAdmin,
			"username": cases.Title(language.English).String(user.Username),
			"date":     utils.GetCurrentTimeFormatted(),
		}

		templateSite := template.Must(
			template.ParseFS(
				env.Files,
				"main/frontend/public/dashboard/index.gohtml",
				"main/frontend/templates/base.gohtml"))

		err = templateSite.Execute(c.Writer, data)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
			return
		}
	})

	r.GET("/bookmanager", middleware.JWTAuth(), func(c *gin.Context) {
		templateSite := template.Must(
			template.ParseFS(
				env.Files,
				"main/frontend/public/bookmanager/index.gohtml",
				"main/frontend/templates/base.gohtml"))

		err := templateSite.Execute(c.Writer, nil)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
			return
		}
	})

	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("auth_token", "", -1, "/", "", false, false)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	return true

	// TODO: add /privacy
	// TODO: add /terms
}

func InitBackendRoutes(r *gin.Engine) bool {
	apiV1 := r.Group("/api/v1")
	v1.ApiHandler(apiV1)

	return true
}

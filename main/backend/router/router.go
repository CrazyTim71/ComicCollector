package router

import (
	v1 "ComicCollector/main/backend/api/v1"
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/middleware"
	"ComicCollector/main/backend/utils/env"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	r.GET("/dashboard", middleware.CheckJwtToken(), func(c *gin.Context) {
		// get the userId
		// because of middleware.CheckJwtToken() we can safely assume that the user id logged in
		id, exists := c.Get("userId")
		if !exists {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		userId, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
			return
		}

		// check if the user is an admin
		isAdmin, err := groups.CheckUserGroup(userId, groups.Administrator)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
			return
		}

		data := map[string]interface{}{
			"isAdmin": isAdmin,
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

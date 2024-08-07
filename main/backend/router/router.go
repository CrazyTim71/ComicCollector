package router

import (
	v1 "ComicCollector/main/backend/api/v1"
	"ComicCollector/main/backend/middleware"
	"ComicCollector/main/backend/utils/env"
	"github.com/gin-gonic/gin"
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

		templateSite := template.Must(
			template.ParseFS(
				env.Files,
				"main/frontend/public/index.gohtml",
				"main/frontend/templates/base.gohtml"))

		err = templateSite.Execute(c.Writer, nil)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
		}
	})

	r.GET("/login", func(c *gin.Context) {
		authCookie, err := c.Cookie("auth_token")
		if authCookie != "" && err == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
			return
		}

		templateSite := template.Must(
			template.ParseFS(
				env.Files,
				"main/frontend/public/login/index.gohtml",
				"main/frontend/templates/base.gohtml"))

		err = templateSite.Execute(c.Writer, nil)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
		}
	})

	r.GET("/dashboard", middleware.CheckJwtToken(), func(c *gin.Context) {
		templateSite := template.Must(
			template.ParseFS(
				env.Files,
				"main/frontend/public/dashboard/index.gohtml",
				"main/frontend/templates/base.gohtml"))

		err := templateSite.Execute(c.Writer, nil)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
		}
	})

	return true
}

func InitBackendRoutes(r *gin.Engine) bool {
	apiV1 := r.Group("/api/v1")
	v1.ApiHandler(apiV1)

	return true
}

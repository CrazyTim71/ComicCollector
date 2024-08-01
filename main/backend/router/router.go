package router

import (
	v1 "ComicCollector/main/backend/api/v1"
	"ComicCollector/main/backend/middleware"
	"ComicCollector/main/backend/utils/env"
	"github.com/gin-gonic/gin"
	"html/template"
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

		template := template.Must(
			template.ParseFS(
				env.Files,
				"main/frontend/public/index.gohtml",
				"main/frontend/templates/base.gohtml"))

		template.Execute(c.Writer, nil)
	})

	r.GET("/login", func(c *gin.Context) {
		authCookie, err := c.Cookie("auth_token")
		if authCookie != "" && err == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
			return
		}

		template := template.Must(
			template.ParseFS(
				env.Files,
				"main/frontend/public/login/index.gohtml",
				"main/frontend/templates/base.gohtml"))

		template.Execute(c.Writer, nil)
	})

	r.GET("/dashboard", middleware.CheckJwtToken(), func(c *gin.Context) {
		template := template.Must(
			template.ParseFS(
				env.Files,
				"main/frontend/public/dashboard/index.gohtml",
				"main/frontend/templates/base.gohtml"))

		template.Execute(c.Writer, nil)
	})

	return true
}

func InitBackendRoutes(r *gin.Engine) bool {
	apiV1 := r.Group("/api/v1")
	v1.ApiHandler(apiV1)

	return true
}

package router

import (
	v1 "ComicCollector/main/backend/api/v1"
	"ComicCollector/main/backend/utils/env"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func InitRouter(r *gin.Engine) bool {
	apiV1 := r.Group("/api/v1")
	v1.ApiHandler(apiV1)

	r.StaticFileFS("/favicon.ico", "main/frontend/static/favicon.ico", http.FS(env.Files))

	r.GET("/login", func(c *gin.Context) {
		template := template.Must(
			template.ParseFS(
				env.Files,
				"main/frontend/public/login/index.gohtml",
				"main/frontend/templates/base.gohtml"))

		template.Execute(c.Writer, nil)
	})

	return true
}

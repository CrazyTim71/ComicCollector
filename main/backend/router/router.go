package router

import (
	v1 "ComicCollector/backend/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) bool {
	apiV1 := r.Group("/api/v1")
	v1.ApiHandler(apiV1)

	return true
}

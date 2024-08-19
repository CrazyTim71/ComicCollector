package v1

import (
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/middleware"
	"github.com/gin-gonic/gin"
)

func BookEditionHandler(rg *gin.RouterGroup) {
	rg.GET("",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		func(c *gin.Context) {
			//bookEditions, err := operations.GetAll(database.MongoDB)
			//if err != nil {
			//	log.Println(err)
			//	c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
			//	return
			//}
			//
			//if bookEditions == nil {
			//	bookEditions = []models.BookEdition{}
			//}
			//
			//c.JSON(http.StatusOK, bookEditions)
		})
}

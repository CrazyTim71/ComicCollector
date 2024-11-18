package router

import (
	v1 "ComicCollector/main/backend/api/v1"
	"ComicCollector/main/backend/utils/env"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
)

var uiFiles, _ = fs.Sub(env.FrontendFiles, "main/frontend/dist")

func InitFrontendRoutes(r *gin.Engine) bool {
	// https://www.reddit.com/r/golang/comments/15vxaxq/q_integrating_vuejs_embedfs_gin/
	// https://github.com/zincsearch/zincsearch/blob/0652db6d39badc753f28ee1122dcbc0e5da1c35e/pkg/routes/routes.go#L43

	// Serve static files from the embedded file system
	// https://stackoverflow.com/questions/62293398/cant-serve-vue-js-spa-app-using-the-noroute-function
	// https://github.com/gin-contrib/static/issues/19#issuecomment-1949903562
	// https://github.com/gin-contrib/static/blob/21b6603afc68fb94b5c7959764f9c198eb9cab52/_example/embed/example.go#L1-L27
	// r.Use(static.Serve("/", static.LocalFile("main/frontend/dist", true)))
	r.Use(static.Serve("/", static.EmbedFolder(env.FrontendFiles, "main/frontend/dist")))
	r.NoRoute(func(c *gin.Context) {
		// This will serve the index.html file from the embedded file system
		// That way the Vue.js app will be able to handle the routing for the frontend
		// c.File("main/frontend/dist/index.html")
		// https://github.com/gin-gonic/gin/issues/2654#issuecomment-944051505
		c.FileFromFS("main/frontend/dist/", http.FS(env.FrontendFiles))
	})

	//r.GET("/", func(c *gin.Context) {
	//	authCookie, err := c.Cookie("auth_token")
	//	if authCookie != "" && err == nil {
	//		c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
	//		return
	//	}
	//})
	//
	//r.GET("/login", func(c *gin.Context) {
	//	authCookie, err := c.Cookie("auth_token")
	//	if authCookie != "" && err == nil {
	//		c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
	//		return
	//	}
	//
	//	SignupEnabled := env.GetSignupEnabled()
	//	data := map[string]interface{}{
	//		"SIGNUP_ENABLED": SignupEnabled,
	//	}
	//
	//	templateSite := template.Must(
	//		template.ParseFS(
	//			env.Files,
	//			"main/frontend/public/login/index.gohtml",
	//			"main/frontend/templates/base.gohtml"))
	//
	//	err = templateSite.Execute(c.Writer, data)
	//	if err != nil {
	//		log.Println(err)
	//		c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
	//		return
	//	}
	//})
	//
	//r.GET("/register", func(c *gin.Context) {
	//	authCookie, err := c.Cookie("auth_token")
	//	if authCookie != "" && err == nil {
	//		c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
	//		return
	//	}
	//
	//	if !env.GetSignupEnabled() {
	//		c.Redirect(http.StatusTemporaryRedirect, "/login")
	//		return
	//	}
	//
	//	templateSite := template.Must(
	//		template.ParseFS(
	//			env.Files,
	//			"main/frontend/public/register/index.gohtml",
	//			"main/frontend/templates/base.gohtml"))
	//
	//	err = templateSite.Execute(c.Writer, nil)
	//	if err != nil {
	//		log.Println(err)
	//		c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
	//		return
	//	}
	//})
	//
	//r.GET("/dashboard", middleware.JWTAuth(), func(c *gin.Context) {
	//	// get the userId
	//	// because of middleware.JWTAuth() we can safely assume that the user id logged in
	//	userId, err := utils.GetUserId(c)
	//	if err != nil {
	//		log.Println(err)
	//		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
	//		return
	//	}
	//
	//	// check if the user is an admin
	//	isAdmin, err := groups.CheckUserGroup(userId, groups.Administrator)
	//	if err != nil {
	//		log.Println(err)
	//		c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
	//		return
	//	}
	//
	//	// get the username
	//	user, err := operations.GetOneById[models.User](database.Tables.User, userId)
	//	if err != nil {
	//		log.Println(err)
	//		c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
	//		return
	//	}
	//
	//	data := map[string]interface{}{
	//		"isAdmin":  isAdmin,
	//		"username": cases.Title(language.English).String(user.Username),
	//		"date":     utils.GetCurrentTimeFormatted(),
	//	}
	//
	//	templateSite := template.Must(
	//		template.ParseFS(
	//			env.Files,
	//			"main/frontend/public/dashboard/index.gohtml",
	//			"main/frontend/templates/base.gohtml"))
	//
	//	err = templateSite.Execute(c.Writer, data)
	//	if err != nil {
	//		log.Println(err)
	//		c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
	//		return
	//	}
	//})
	//
	//r.GET("/bookmanager", middleware.JWTAuth(), func(c *gin.Context) {
	//	templateSite := template.Must(
	//		template.ParseFS(
	//			env.Files,
	//			"main/frontend/public/bookmanager/index.gohtml",
	//			"main/frontend/templates/base.gohtml"))
	//
	//	err := templateSite.Execute(c.Writer, nil)
	//	if err != nil {
	//		log.Println(err)
	//		c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred while rendering the templateSite", "error": true})
	//		return
	//	}
	//})
	//
	//r.GET("/logout", func(c *gin.Context) {
	//	c.SetCookie("auth_token", "", -1, "/", "", false, false)
	//	c.Redirect(http.StatusTemporaryRedirect, "/")
	//})
	//
	return true

	// TODO: add /privacy
	// TODO: add /terms
}

func InitBackendRoutes(r *gin.Engine) bool {
	// TODO: add CORS

	apiV1 := r.Group("/api/v1")
	v1.ApiHandler(apiV1)

	return true
}

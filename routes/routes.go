package routes

import (
	"kwickiwiki/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./static")

	// Routes for viewing, editing, and searching wiki pages
	router.GET("/", func(c *gin.Context) { c.Redirect(302, "/home") })
	router.GET("/:title", handlers.AuthMiddleware(), handlers.ShowPage)
	router.GET("/:title/edit", handlers.AuthMiddleware(), handlers.EditPage)
	router.POST("/:title/save", handlers.AuthMiddleware(), handlers.SavePage)
	router.GET("/search", handlers.SearchPage)
	router.GET("/register", handlers.ShowRegisterPage)
	router.POST("/login", handlers.LoginUser)

	return router
}

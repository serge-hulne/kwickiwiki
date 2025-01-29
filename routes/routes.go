package routes

import (
	"wiki_project/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./static")

	// Routes for viewing, editing, and searching wiki pages
	router.GET("/", func(c *gin.Context) { c.Redirect(302, "/home") })
	router.GET("/:title", handlers.ShowPage)
	router.GET("/:title/edit", handlers.EditPage)
	router.POST("/:title/save", handlers.SavePage)
	router.GET("/search", handlers.SearchPage) // ✅ Add search route

	return router
}

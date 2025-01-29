package main

import (
	"kwickiwiki/handlers"
	"kwickiwiki/models"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	models.InitDB()

	// Create a new Gin router
	r := gin.Default()

	// Initialize session store (cookie-based)
	store := cookie.NewStore([]byte("your-secret-key"))
	r.Use(sessions.Sessions("session_name", store))

	// Define routes and apply middleware
	r.GET("/", handlers.ShowPage)

	r.GET("/register", handlers.ShowRegisterPage)
	r.POST("/register", handlers.RegisterUser)

	r.GET("/login", handlers.LoginUser)
	r.POST("/login", handlers.LoginUser)

	// Protect these routes with authentication
	r.GET("/edit/:title", handlers.AuthMiddleware(), handlers.EditPage)
	r.POST("/edit/:title", handlers.AuthMiddleware(), handlers.SavePage)

	// Run the server
	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}

// Assuming the home page is just a wiki page

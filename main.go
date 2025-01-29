// package main

// import (
// 	"wiki_project/models"
// 	"wiki_project/routes"

// 	"github.com/gin-contrib/sessions"
// 	"github.com/gin-contrib/sessions/cookie"
// )

// func main() {
// 	models.InitDB() // Initialize the database

// 	// Set up session middleware
// 	store := cookie.NewStore([]byte("your-secret-key"))

// 	router := routes.SetupRouter()

// 	router.Use(sessions.Sessions("session_name", store))

// 	router.Run(":8080")
// }

package main

import (
	"log"

	"wiki_project/handlers"
	"wiki_project/models"

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
	r.GET("/", handlers.HomePage)
	r.GET("/register", handlers.ShowRegisterPage)
	r.POST("/register", handlers.RegisterUser)
	r.GET("/login", handlers.ShowLoginPage)
	r.POST("/login", handlers.LoginUser)

	// Protect these routes with authentication
	r.GET("/edit/:title", handlers.AuthMiddleware(), handlers.ShowEditPage)
	r.POST("/edit/:title", handlers.AuthMiddleware(), handlers.SaveEditedPage)

	// Run the server
	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}

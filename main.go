package main

import (
	"wiki_project/models"
	"wiki_project/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func main() {
	models.InitDB() // Initialize the database

	// Set up session middleware
	store := cookie.NewStore([]byte("your-secret-key"))

	router := routes.SetupRouter()

	router.Use(sessions.Sessions("session_name", store))

	router.Run(":8080")
}

package main

import (
	"wiki_project/models"
	"wiki_project/routes"
)

func main() {
	models.InitDB() // Initialize the database
	router := routes.SetupRouter()
	router.Run(":8080")
}

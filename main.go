package main

import (
	"github.com/Stand/db"
	"github.com/Stand/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize DB
	db.InitDB()

	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}

package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/vehicles", getVehicles)
	server.GET("/vehicles/:id", getVehicle)
	server.POST("/vehicles", createVehicle)
	server.PUT("/vehicles/:id", updateVehicle)
	server.DELETE("/vehicles/:id", deleteVehicle)
	// /vehicles?type=carro&brand=Toyota&year=2020

	server.GET("/clients", getClients)
	server.GET("/clients/:id", getClient)
	server.POST("/clients", createClient)
	server.PUT("/clients/:id", updateClient)
	server.DELETE("/clients/:id", deleteClient)

	server.GET("/sales", getSales)
	server.GET("/sales/:id", getSale)
	server.POST("/sales", createSale)
}

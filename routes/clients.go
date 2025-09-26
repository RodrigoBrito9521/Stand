package routes

import (
	"net/http"
	"strconv"

	"github.com/Stand/models"
	"github.com/gin-gonic/gin"
)

func getClients(context *gin.Context) {
	clients, err := models.GetAllClients()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch clients. Try again later."})
		return
	}
	context.JSON(http.StatusOK, clients)
}

func getClient(context *gin.Context) {
	clientId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse client id."})
		return
	}

	client, err := models.GetClientByID(clientId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch client."})
		return
	}

	context.JSON(http.StatusOK, client)
}

func createClient(context *gin.Context) {
	var client models.Client
	err := context.ShouldBindJSON(&client)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = client.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create client. Try again later."})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Client created!", "client": client})
}

func updateClient(context *gin.Context) {
	clientId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse client id."})
		return
	}

	_, err = models.GetClientByID(clientId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the client."})
		return
	}

	// If you need authentication, implement it separately

	var updatedClient models.Client
	err = context.ShouldBindJSON(&updatedClient)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	updatedClient.ID = clientId
	err = updatedClient.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update client."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Client updated successfully!"})
}

func deleteClient(context *gin.Context) {
	clientId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse client id."})
		return
	}

	client, err := models.GetClientByID(clientId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the client."})
		return
	}

	// If you need authentication, implement it as middleware

	err = client.Delete()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the client."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully!"})
}

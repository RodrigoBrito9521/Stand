package routes

import (
	"net/http"
	"strconv"

	"github.com/Stand/models"
	"github.com/gin-gonic/gin"
)

func getVehicles(context *gin.Context) {
	vehicle, err := models.GetAllVehicles()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch vehicles. Try again later."})
		return
	}
	context.JSON(http.StatusOK, vehicle)
}

func getVehicle(context *gin.Context) {
	vehicleId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse vehicle id."})
		return
	}

	vehicle, err := models.GetVehicleByID(vehicleId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch vehicle."})
		return
	}

	context.JSON(http.StatusOK, vehicle)
}

func createVehicle(context *gin.Context) {
	var vehicle models.Vehicle
	err := context.ShouldBindJSON(&vehicle)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = vehicle.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create vehicle. Try again later."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Vehicle created!", "vehicle": vehicle})

}

func updateVehicle(context *gin.Context) {
	vehicleId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse vehicle id."})
		return
	}

	var updateVehicle models.Vehicle
	err = context.ShouldBindJSON(&updateVehicle)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	updateVehicle.ID = int(vehicleId)
	err = updateVehicle.UpdateVehicle()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update vehicle."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Vehicle updated successfully!"})

}

func deleteVehicle(context *gin.Context) {
	vehicleID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse vehicle id."})
		return
	}

	vehicle, err := models.GetVehicleByID(vehicleID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the vehicle."})
		return
	}

	err = vehicle.Delete()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the vehicle."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Vehicle deleted successfully!"})
}

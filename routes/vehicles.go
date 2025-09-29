package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Stand/models"
	"github.com/gin-gonic/gin"
)

func getVehicles(context *gin.Context) {
	vehicleType := context.Query("type")
	brand := context.Query("brand")
	yearStr := context.Query("year")

	var year int
	if yearStr != "" {
		parsedYear, err := strconv.Atoi(yearStr)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid year format."})
			return
		}
		year = parsedYear
	}

	// If any filters are provided, use filtered search
	if vehicleType != "" || brand != "" || year > 0 {
		vehicles, err := models.GetVehiclesWithFilters(vehicleType, brand, year)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch vehicles. Try again later."})
			return
		}
		context.JSON(http.StatusOK, vehicles)
		return
	}

	// Otherwise, get all vehicles
	vehicles, err := models.GetAllVehicles()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch vehicles. Try again later."})
		return
	}
	context.JSON(http.StatusOK, vehicles)
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
	log.Println("Starting createVehicle handler")

	var vehicle models.Vehicle
	err := context.ShouldBindJSON(&vehicle)

	if err != nil {
		log.Printf("JSON binding error: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	log.Printf("Parsed vehicle data: %+v", vehicle)

	err = vehicle.Save()

	if err != nil {
		log.Printf("Database save error: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create vehicle. Try again later."})
		return
	}

	log.Println("Vehicle created successfully")
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

package routes

import (
	"net/http"
	"strconv"

	"github.com/Stand/models"
	"github.com/gin-gonic/gin"
)

func getSales(context *gin.Context) {
	sales, err := models.GetAllSales()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch sales. Try again later."})
		return
	}
	context.JSON(http.StatusOK, sales)
}

func getSale(context *gin.Context) {
	saleId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse sale id."})
		return
	}

	sale, err := models.GetSaleByID(saleId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch sale."})
		return
	}

	context.JSON(http.StatusOK, sale)
}

func createSale(context *gin.Context) {
	var sale models.Sale
	err := context.ShouldBindJSON(&sale)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = sale.Save()

	if err != nil {
		// Check if it's a vehicle already sold error
		if vehicleErr, ok := err.(*models.VehicleAlreadySoldError); ok {
			context.JSON(http.StatusConflict, gin.H{
				"message":    "Vehicle is already sold",
				"vehicle_id": vehicleErr.VehicleID,
			})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create sale. Try again later."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Sale created successfully!", "sale": sale})
}

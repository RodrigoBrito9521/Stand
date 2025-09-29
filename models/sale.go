package models

import (
	"database/sql"
	"log"
	"time"

	"github.com/Stand/db"
)

type Sale struct {
	ID        int64     `json:"id"`
	ClientID  int64     `json:"client_id" binding:"required"`
	VehicleID int64     `json:"vehicle_id" binding:"required"`
	Price     float64   `json:"price" binding:"required"`
	SaleDate  time.Time `json:"sale_date"`
}

// SaleWithDetails represents a sale with client and vehicle information
type SaleWithDetails struct {
	ID       int64     `json:"id"`
	Price    float64   `json:"price"`
	SaleDate time.Time `json:"sale_date"`
	Client   Client    `json:"client"`
	Vehicle  Vehicle   `json:"vehicle"`
}

func (s *Sale) Save() error {
	log.Printf("[v0] Starting Sale.Save() with data: %+v", s)

	// First check if vehicle is already sold
	var existingSaleID int64
	checkQuery := "SELECT id FROM sales WHERE vehicle_id = $1"
	err := db.DB.QueryRow(checkQuery, s.VehicleID).Scan(&existingSaleID)

	if err != sql.ErrNoRows {
		if err == nil {
			log.Printf("[v0] Vehicle %d is already sold", s.VehicleID)
			return &VehicleAlreadySoldError{VehicleID: s.VehicleID}
		}
		log.Printf("[v0] Error checking existing sale: %v", err)
		return err
	}

	// Check if vehicle exists and is available
	vehicle, err := GetVehicleByID(s.VehicleID)
	if err != nil {
		log.Printf("[v0] Error getting vehicle: %v", err)
		return err
	}

	if vehicle.Status == "sold" {
		log.Printf("[v0] Vehicle %d status is already 'sold'", s.VehicleID)
		return &VehicleAlreadySoldError{VehicleID: s.VehicleID}
	}

	// Check if client exists
	_, err = GetClientByID(s.ClientID)
	if err != nil {
		log.Printf("[v0] Error getting client: %v", err)
		return err
	}

	// Create the sale
	s.SaleDate = time.Now()
	query := `INSERT INTO sales (client_id, vehicle_id, price, sale_date)
	VALUES ($1, $2, $3, $4) RETURNING id`

	log.Printf("[v0] SQL Query: %s", query)
	log.Printf("[v0] Parameters: client_id=%d, vehicle_id=%d, price=%.2f, sale_date=%v",
		s.ClientID, s.VehicleID, s.Price, s.SaleDate)

	err = db.DB.QueryRow(query, s.ClientID, s.VehicleID, s.Price, s.SaleDate).Scan(&s.ID)
	if err != nil {
		log.Printf("[v0] QueryRow/Scan error: %v", err)
		return err
	}

	log.Printf("[v0] Sale saved successfully with ID: %d", s.ID)

	// Update vehicle status to sold
	updateQuery := "UPDATE vehicles SET status = 'sold' WHERE id = $1"
	_, err = db.DB.Exec(updateQuery, s.VehicleID)
	if err != nil {
		log.Printf("[v0] Error updating vehicle status: %v", err)
		return err
	}

	log.Printf("[v0] Vehicle %d status updated to 'sold'", s.VehicleID)
	return nil
}

func GetAllSales() ([]SaleWithDetails, error) {
	query := `
	SELECT 
		s.id, s.price, s.sale_date,
		c.id, c.name, c.email, c.phone,
		v.id, v.type, v.brand, v.model, v.year, v.motor, v.status
	FROM sales s
	JOIN clients c ON s.client_id = c.id
	JOIN vehicles v ON s.vehicle_id = v.id
	ORDER BY s.sale_date DESC`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []SaleWithDetails

	for rows.Next() {
		var sale SaleWithDetails
		err := rows.Scan(
			&sale.ID, &sale.Price, &sale.SaleDate,
			&sale.Client.ID, &sale.Client.Name, &sale.Client.Email, &sale.Client.Phone,
			&sale.Vehicle.ID, &sale.Vehicle.Type, &sale.Vehicle.Brand,
			&sale.Vehicle.Model, &sale.Vehicle.Year, &sale.Vehicle.Motor, &sale.Vehicle.Status,
		)
		if err != nil {
			return nil, err
		}
		sales = append(sales, sale)
	}

	return sales, nil
}

func GetSaleByID(id int64) (*SaleWithDetails, error) {
	query := `
	SELECT 
		s.id, s.price, s.sale_date,
		c.id, c.name, c.email, c.phone,
		v.id, v.type, v.brand, v.model, v.year, v.motor, v.status
	FROM sales s
	JOIN clients c ON s.client_id = c.id
	JOIN vehicles v ON s.vehicle_id = v.id
	WHERE s.id = $1`

	row := db.DB.QueryRow(query, id)

	var sale SaleWithDetails
	err := row.Scan(
		&sale.ID, &sale.Price, &sale.SaleDate,
		&sale.Client.ID, &sale.Client.Name, &sale.Client.Email, &sale.Client.Phone,
		&sale.Vehicle.ID, &sale.Vehicle.Type, &sale.Vehicle.Brand,
		&sale.Vehicle.Model, &sale.Vehicle.Year, &sale.Vehicle.Motor, &sale.Vehicle.Status,
	)

	if err != nil {
		return nil, err
	}

	return &sale, nil
}

// Custom error for vehicle already sold
type VehicleAlreadySoldError struct {
	VehicleID int64
}

func (e *VehicleAlreadySoldError) Error() string {
	return "vehicle is already sold"
}

/*
// Helper functions to get client and vehicle by ID (standalone functions)

	func GetClientByID(id int64) (*Client, error) {
		query := "SELECT id, name, email, phone FROM clients WHERE id=$1"
		row := db.DB.QueryRow(query, id)

		var client Client
		err := row.Scan(&client.ID, &client.Name, &client.Email, &client.Phone)
		if err != nil {
			return nil, err
		}

		return &client, nil
	}

	func GetVehicleByID(id int64) (*Vehicle, error) {
		query := "SELECT id, type, brand, model, year, motor, status FROM vehicles WHERE id=$1"
		row := db.DB.QueryRow(query, id)

		var vehicle Vehicle
		err := row.Scan(&vehicle.ID, &vehicle.Type, &vehicle.Brand, &vehicle.Model, &vehicle.Year, &vehicle.Motor, &vehicle.Status)
		if err != nil {
			return nil, err
		}

		return &vehicle, nil
	}

type Client struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Vehicle struct {
	ID     int64  `json:"id"`
	Type   string `json:"type"`
	Brand  string `json:"brand"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	Motor  string `json:"motor"`
	Status string `json:"status"`
}
*/

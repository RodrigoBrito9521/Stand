package models

import (
	"database/sql"
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

type SaleWithDetails struct {
	ID       int64     `json:"id"`
	Price    float64   `json:"price"`
	SaleDate time.Time `json:"sale_date"`
	Client   Client    `json:"client"`
	Vehicle  Vehicle   `json:"vehicle"`
}

func (s *Sale) Save() error {
	// ve se o vehicle ja foi vendido na tabela sales
	var existingSaleID int64
	checkQuery := "SELECT id FROM sales WHERE vehicle_id = ?"
	err := db.DB.QueryRow(checkQuery, s.VehicleID).Scan(&existingSaleID)

	if err != sql.ErrNoRows {
		if err == nil {
			return &VehicleAlreadySoldError{VehicleID: s.VehicleID}
		}
		return err
	}

	// Ve se o vehicle existe e se não foi já registado como sold
	vehicle, err := GetVehicleByID(s.VehicleID)
	if err != nil {
		return err
	}

	if vehicle.Status == "sold" {
		return &VehicleAlreadySoldError{VehicleID: s.VehicleID}
	}

	// Ve se o cliente existe
	_, err = GetClientByID(s.ClientID)
	if err != nil {
		return err
	}

	// Create the sale
	s.SaleDate = time.Now()
	query := `INSERT INTO sales (client_id, vehicle_id, price, sale_date)
	VALUES (?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(s.ClientID, s.VehicleID, s.Price, s.SaleDate)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	s.ID = id

	// Update vehicle status to sold
	updateQuery := "UPDATE vehicles SET status = 'sold' WHERE id = ?"
	_, err = db.DB.Exec(updateQuery, s.VehicleID)

	return err
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
	WHERE s.id = ?`

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

// Helper functions to get client and vehicle by ID (standalone functions)
/*
func GetClientByID(id int64) (*Client, error) {
	query := "SELECT id, name, email, phone FROM clients WHERE id=?"
	row := db.DB.QueryRow(query, id)

	var client Client
	err := row.Scan(&client.ID, &client.Name, &client.Email, &client.Phone)
	if err != nil {
		return nil, err
	}

	return &client, nil

*/

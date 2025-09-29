package models

import (
	"fmt"
	"github.com/Stand/db"
	"log"
)

type Vehicle struct {
	ID     int
	Type   string `binding:"required"`
	Brand  string `binding:"required"`
	Model  string `binding:"required"`
	Year   int    `binding:"required"`
	Motor  string `binding:"required"`
	Status string `binding:"required"`
}

var vehicles = []Vehicle{}

func (v *Vehicle) Save() error {
	log.Printf("[v0] Starting Vehicle.Save() with data: %+v", v)

	query := `
	INSERT INTO vehicles(type, brand, model, year, motor, status)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	log.Printf("[v0] SQL Query: %s", query)
	log.Printf("[v0] Parameters: type=%s, brand=%s, model=%s, year=%d, motor=%s, status=%s",
		v.Type, v.Brand, v.Model, v.Year, v.Motor, v.Status)

	err := db.DB.QueryRow(query, v.Type, v.Brand, v.Model, v.Year, v.Motor, v.Status).Scan(&v.ID)
	if err != nil {
		log.Printf("[v0] QueryRow/Scan error: %v", err)
		return err
	}

	log.Printf("[v0] Vehicle saved successfully with ID: %d", v.ID)
	return nil
}

func GetAllVehicles() ([]Vehicle, error) {
	query := "SELECT * FROM vehicles"
	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var vehicles []Vehicle

	for rows.Next() {
		var vehicle Vehicle
		err := rows.Scan(&vehicle.ID, &vehicle.Type, &vehicle.Brand, &vehicle.Model, &vehicle.Year, &vehicle.Motor, &vehicle.Status)

		if err != nil {
			return nil, err
		}

		vehicles = append(vehicles, vehicle)
	}
	return vehicles, nil
}

func GetVehicleByID(id int64) (*Vehicle, error) {
	query := "SELECT * FROM vehicles WHERE id=$1"
	row := db.DB.QueryRow(query, id)

	var vehicle Vehicle
	err := row.Scan(&vehicle.ID, &vehicle.Type, &vehicle.Brand, &vehicle.Model, &vehicle.Year, &vehicle.Motor, &vehicle.Status)

	if err != nil {
		return nil, err
	}

	return &vehicle, nil
}

func (vehicle Vehicle) UpdateVehicle() error {
	query := `UPDATE vehicles 
	SET type=$1, brand=$2, model=$3, year=$4, motor=$5, status=$6
	WHERE id=$7
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(vehicle.Type, vehicle.Brand, vehicle.Model, vehicle.Year, vehicle.Motor, vehicle.Status, vehicle.ID)
	return err
}

func (vehicle Vehicle) Delete() error {
	query := "DELETE FROM vehicles WHERE id = $1"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(vehicle.ID)
	return err
}

func GetVehiclesWithFilters(vehicleType, brand string, year int) ([]Vehicle, error) {
	query := "SELECT * FROM vehicles WHERE 1=1"
	var args []interface{}
	paramCount := 1

	if vehicleType != "" {
		query += " AND type = $" + fmt.Sprintf("%d", paramCount)
		args = append(args, vehicleType)
		paramCount++
	}

	if brand != "" {
		query += " AND brand = $" + fmt.Sprintf("%d", paramCount)
		args = append(args, brand)
		paramCount++
	}

	if year > 0 {
		query += " AND year = $" + fmt.Sprintf("%d", paramCount)
		args = append(args, year)
		paramCount++
	}

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var vehicle Vehicle
		err := rows.Scan(&vehicle.ID, &vehicle.Type, &vehicle.Brand, &vehicle.Model, &vehicle.Year, &vehicle.Motor, &vehicle.Status)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

package models

import (
	"github.com/Stand/db"
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
	query := `
	INSERT INTO vehicles(type, brand, model, year, motor, status)
	VALUES (?, ?, ?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(v.Type, v.Brand, v.Model, v.Year, v.Motor, v.Status)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	v.ID = int(id)

	return err
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
	query := "SELECT * FROM vehicles WHERE id=?"
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
	SET type=?, brand=?, model=?, year=?, motor=?, status=?
	WHERE id=?
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
	query := "DELETE FROM vehicles WHERE id = ?"
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

	if vehicleType != "" {
		query += " AND type = ?"
		args = append(args, vehicleType)
	}

	if brand != "" {
		query += " AND brand = ?"
		args = append(args, brand)
	}

	if year > 0 {
		query += " AND year = ?"
		args = append(args, year)
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

package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	connStr := "host=db user=postgres password=postgres dbname=stand_automovel port=5432 sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao conectar à base de dados:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Erro ao fazer ping à base de dados:", err)
	}

	fmt.Println("Conexão à base de dados PostgreSQL estabelecida com sucesso!")

	createTables()
}

func createTables() {
	createVehiclesTable := `
    CREATE TABLE IF NOT EXISTS vehicles (
        id SERIAL PRIMARY KEY,
        type TEXT NOT NULL,
        brand TEXT NOT NULL,
        model TEXT NOT NULL,
        year INTEGER NOT NULL,
        motor TEXT NOT NULL,
        status TEXT NOT NULL
    )`

	_, err := DB.Exec(createVehiclesTable)
	if err != nil {
		panic("Could not create vehicle table: " + err.Error())
	}

	createClientTable := `
    CREATE TABLE IF NOT EXISTS clients (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT NOT NULL,
        phone BIGINT NOT NULL
    )`

	_, err = DB.Exec(createClientTable)
	if err != nil {
		panic("Could not create client table: " + err.Error())
	}

	createSalesTable := `
    CREATE TABLE IF NOT EXISTS sales (
        id SERIAL PRIMARY KEY,
        client_id INTEGER NOT NULL REFERENCES clients(id),
        vehicle_id INTEGER NOT NULL REFERENCES vehicles(id) UNIQUE,
        price REAL NOT NULL,
        sale_date TIMESTAMP NOT NULL
    )`

	_, err = DB.Exec(createSalesTable)
	if err != nil {
		panic("Could not create sales table: " + err.Error())
	}
}

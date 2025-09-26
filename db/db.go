package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	// SQLite cria o ficheiro automaticamente
	DB, err = sql.Open("sqlite", "./stand_automovel.db")
	if err != nil {
		log.Fatal("Erro ao conectar à base de dados:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Erro ao fazer ping à base de dados:", err)
	}

	fmt.Println("Conexão à base de dados SQLite estabelecida com sucesso!")

	createTables()

}

func createTables() {
	createVehiclesTable := `
	CREATE TABLE IF NOT EXISTS vehicles (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    type TEXT NOT NULL,
	    brand TEXT NOT NULL,
	    model TEXT NOT NULL,
	    year INTEGER NOT NULL,
	    motor TEXT NOT NULL,
	    status TEXT NOT NULL
	)`

	_, err := DB.Exec(createVehiclesTable)

	if err != nil {
		panic("Could not create vehicle table.")
	}

	createClientTable := `
	CREATE TABLE IF NOT EXISTS clients (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT NOT NULL,
	    email TEXT NOT NULL,
	    phone INTEGER NOT NULL
	)`

	_, err = DB.Exec(createClientTable)

	if err != nil {
		panic("Could not create client table.")
	}

	createSalesTable := `
	CREATE TABLE IF NOT EXISTS sales (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    client_id INTEGER NOT NULL,
	    vehicle_id INTEGER NOT NULL,
	    price REAL NOT NULL,
	    sale_date DATETIME NOT NULL,
	    FOREIGN KEY (client_id) REFERENCES clients(id),
	    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
	    UNIQUE(vehicle_id)
	)`

	_, err = DB.Exec(createSalesTable)

	if err != nil {
		panic("Could not create sales table.")
	}
}

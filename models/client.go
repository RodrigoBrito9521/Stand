package models

import (
	"github.com/Stand/db"
	"log"
)

type Client struct {
	ID    int64  `json:"id"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Phone int    `json:"phone" binding:"required"`
}

func (c *Client) Save() error {
	log.Printf("[v0] Starting Client.Save() with data: %+v", c)

	query := `INSERT INTO clients (name, email, phone)
	VALUES ($1, $2, $3) RETURNING id`

	log.Printf("[v0] SQL Query: %s", query)
	log.Printf("[v0] Parameters: name=%s, email=%s, phone=%d", c.Name, c.Email, c.Phone)

	err := db.DB.QueryRow(query, c.Name, c.Email, c.Phone).Scan(&c.ID)
	if err != nil {
		log.Printf("[v0] QueryRow/Scan error: %v", err)
		return err
	}

	log.Printf("[v0] Client saved successfully with ID: %d", c.ID)
	return nil
}

func GetAllClients() ([]Client, error) {
	query := "SELECT * FROM clients"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []Client

	for rows.Next() {
		var client Client
		err := rows.Scan(&client.ID, &client.Name, &client.Email, &client.Phone)
		if err != nil {
			return clients, err
		}
		clients = append(clients, client)
	}
	return clients, err
}

func GetClientByID(id int64) (*Client, error) {
	query := "SELECT * FROM clients WHERE id=$1"
	row := db.DB.QueryRow(query, id)

	var client Client
	err := row.Scan(&client.ID, &client.Name, &client.Email, &client.Phone)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (c *Client) Update() error {
	query := `
	UPDATE clients
	SET name=$1,email=$2,phone=$3
	WHERE id=$4`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		log.Printf("[v0] Prepare error: %v", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(c.Name, c.Email, c.Phone, c.ID)
	if err != nil {
		log.Printf("[v0] Exec error: %v", err)
		return err
	}

	log.Printf("[v0] Client updated successfully with ID: %d", c.ID)
	return nil
}

func (c *Client) Delete() error {
	query := "DELETE FROM clients WHERE id=$1"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("[v0] Prepare error: %v", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.ID)
	if err != nil {
		log.Printf("[v0] Exec error: %v", err)
		return err
	}

	log.Printf("[v0] Client deleted successfully with ID: %d", c.ID)
	return nil
}

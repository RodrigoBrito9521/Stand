package models

import (
	"github.com/Stand/db"
)

type Client struct {
	ID    int64  `json:"id"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Phone int    `json:"phone" binding:"required"`
}

func (c *Client) Save() error {
	query := `INSERT INTO clients (name, email, phone)
	VALUES (?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(c.Name, c.Email, c.Phone)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	c.ID = id
	return err
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
	query := "SELECT * FROM clients WHERE id=?"
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
	SET name=?,email=?,phone=?
	WHERE id=?`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(c.Name, c.Email, c.Phone, c.ID)
	return err
}

func (c *Client) Delete() error {
	query := "DELETE FROM clients WHERE id=?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.ID)
	return err
}

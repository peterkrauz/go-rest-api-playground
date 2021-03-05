package db

import (
	"database/sql"
	"github.com/peterkrauz/go-rest-api-playground/models"
)

func (db Database) GetAllItems() (*models.ItemList, error) {
	list := &models.ItemList{}
	rows, err := db.Connection.Query("SELECT * FROM items ORDER BY Id desc")
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.Id, &item.Name, &item.Description, &item.CreatedAt)
		if err != nil {
			return list, err
		}
		list.Items = append(list.Items, item)
	}

	return list, nil
}

func (db Database) GetItemById(itemId int) (models.Item, error) {
	item := models.Item{}
	query := `SELECT * FROM items WHERE id = $1`
	row := db.Connection.QueryRow(query, itemId)
	switch err := row.Scan(&item.Id, &item.Name, &item.Description, &item.CreatedAt); err {
	case sql.ErrNoRows:
		return item, ErrorNoItemFound
	default:
		return item, err
	}
}

func (db Database) CreateItem(item *models.Item) error {
	var id int
	var createdAt string

	query := `INSERT INTO items (name, description) VALUES ($1, $2) RETURNING id, created_at`
	err := db.Connection.QueryRow(query, item.Name, item.Description).Scan(&id, &createdAt)
	if err != nil {
		return err
	}

	item.Id = id
	item.CreatedAt = createdAt
	return nil
}

func (db Database) UpdateItem(itemId int, itemData models.Item) (models.Item, error) {
	item := models.Item{}
	query := `UPDATE item SET name=$1, description=$2 WHERE id=$3 RETURNING id, name, description, created_at;`
	err := db.Connection.QueryRow(query, itemData.Name, itemData.Description, itemId).Scan(&item.Id, &item.Name, &item.Description, &item.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return item, ErrorNoItemFound
		}
		return item, err
	}
	return item, nil
}

func (db Database) DeleteItem(itemId int) error {
	query := `DELETE FROM items WHERE id = $1`
	_, err := db.Connection.Exec(query, itemId)
	switch err {
	case sql.ErrNoRows:
		return ErrorNoItemFound
	default:
		return err
	}
}

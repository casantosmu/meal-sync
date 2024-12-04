package models

import (
	"database/sql"
	"strings"
)

type Shopping struct {
	ID          int
	Name        string
	IsPurchased bool
}

type ShoppingModel struct {
	DB *sql.DB
}

func (m ShoppingModel) BulkCreate(names []string) ([]int, error) {
	query := "INSERT INTO shopping (name) VALUES " + strings.Repeat("(?),", len(names))
	query = query[:len(query)-1] // Remove the last comma
	query += " RETURNING shopping_id;"

	args := make([]any, len(names))
	for i, name := range names {
		args[i] = name
	}

	rows, err := m.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int, 0, len(names))
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

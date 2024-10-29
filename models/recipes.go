package models

import (
	"database/sql"
	"errors"
)

type Recipe struct {
	Id          int
	Title       string
	ImageURL    string
	Description string
	Ingredients string
	Directions  string
}

type RecipeModel struct {
	DB *sql.DB
}

func (m RecipeModel) Create(title, img string) (int, error) {
	query := `INSERT INTO recipes (title, img_url)
	VALUES (?, ?)
	RETURNING recipe_id;
	`

	var id int
	err := m.DB.QueryRow(query, title, img).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m RecipeModel) GetByPk(id int) (Recipe, error) {
	query := `SELECT recipe_id, title, img_url, COALESCE(description, ''), COALESCE(ingredients, ''), COALESCE(directions, '')
	FROM recipes
	WHERE recipe_id = ?;`

	r := Recipe{}
	err := m.DB.QueryRow(query, id).Scan(&r.Id, &r.Title, &r.ImageURL, &r.Description, &r.Ingredients, &r.Directions)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return r, ErrNotFound
		}
		return r, err
	}

	return r, nil
}

func (m RecipeModel) Search(search string) ([]Recipe, error) {
	query := "SELECT recipe_id, title, img_url FROM recipes"
	var args []any

	if search != "" {
		query += " WHERE title LIKE ?"
		args = append(args, "%"+search+"%")
	}

	query += " ORDER BY title ASC;"

	rows, err := m.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Recipe
	for rows.Next() {
		var r Recipe
		err := rows.Scan(&r.Id, &r.Title, &r.ImageURL)
		if err != nil {
			return nil, err
		}
		list = append(list, r)
	}

	return list, nil
}

func (m RecipeModel) UpdateByPk(id int, title, description, ingredients, directions string) error {
	query := `UPDATE recipes
	SET title = ?,
		description = ?,
		ingredients = ?,
		directions = ?
	WHERE recipe_id = ?;`

	result, err := m.DB.Exec(query, title, description, ingredients, directions, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

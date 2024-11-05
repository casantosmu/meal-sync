package models

import (
	"database/sql"
	"errors"
	"strings"
)

type Recipe struct {
	ID          int
	Title       string
	ImageURL    string
	Description string
	Ingredients string
	Directions  string
}

func (r Recipe) ImageURLOrDefault() string {
	if r.ImageURL == "" {
		return "/static/images/recipe_placeholder.svg"
	}
	return r.ImageURL
}

func (r Recipe) IngredientsToList() []string {
	lines := strings.Split(r.Ingredients, "\n")
	var ingredients = []string{}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			ingredients = append(ingredients, trimmed)
		}
	}
	if len(ingredients) == 0 {
		return []string{""}
	}
	return ingredients
}

func (r Recipe) DirectionsToList() []string {
	lines := strings.FieldsFunc(r.Directions, func(r rune) bool {
		return r == '\n' || r == '\r'
	})
	var directions = []string{}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			directions = append(directions, trimmed)
		}
	}
	if len(directions) == 0 {
		return []string{""}
	}
	return directions
}

type RecipeModel struct {
	DB *sql.DB
}

func (m RecipeModel) Create(title string) (int, error) {
	query := `INSERT INTO recipes (title)
	VALUES (?)
	RETURNING recipe_id;`

	var id int
	err := m.DB.QueryRow(query, title).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m RecipeModel) GetByPk(id int) (Recipe, error) {
	query := `SELECT recipe_id, title, COALESCE(img_url, ''), COALESCE(description, ''), COALESCE(ingredients, ''), COALESCE(directions, '')
	FROM recipes
	WHERE recipe_id = ?;`

	r := Recipe{}
	err := m.DB.QueryRow(query, id).Scan(&r.ID, &r.Title, &r.ImageURL, &r.Description, &r.Ingredients, &r.Directions)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return r, ErrNotFound
		}
		return r, err
	}

	return r, nil
}

func (m RecipeModel) Search(search string) ([]Recipe, error) {
	query := "SELECT recipe_id, title, COALESCE(img_url, '') FROM recipes"
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
		err := rows.Scan(&r.ID, &r.Title, &r.ImageURL)
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

	result, err := m.DB.Exec(query,
		title,
		newNullString(description),
		newNullString(ingredients),
		newNullString(directions),
		id,
	)
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

func (m RecipeModel) RemoveByPk(id int) error {
	query := "DELETE FROM recipes WHERE recipe_id = ?;"

	result, err := m.DB.Exec(query, id)
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

func (m RecipeModel) UpdateImageByPk(id int, path string) error {
	query := `UPDATE recipes
	SET img_url = ?
	WHERE recipe_id = ?;`

	result, err := m.DB.Exec(query, newNullString(path), id)
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

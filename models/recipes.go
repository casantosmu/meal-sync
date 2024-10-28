package models

import "database/sql"

type Recipe struct {
	Id          int
	Title       string
	Img         string
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
		err := rows.Scan(&r.Id, &r.Title, &r.Img)
		if err != nil {
			return nil, err
		}
		list = append(list, r)
	}

	return list, nil
}

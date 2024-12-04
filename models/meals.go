package models

import (
	"database/sql"
	"time"
)

type Meal struct {
	ID     int
	Date   string
	Recipe Recipe
}

type MealsByDate struct {
	Date  time.Time
	Meals []Meal
}

func (m MealsByDate) DateFormat() string {
	return m.Date.Format(DateFormat)
}

func (m MealsByDate) DayOfWeek() string {
	return m.Date.Weekday().String()
}

func (m MealsByDate) MonthDay() string {
	return m.Date.Format("Jan 2")
}

type MealModel struct {
	DB *sql.DB
}

func (m MealModel) Create(date string, recipeID int) (int, error) {
	query := `INSERT INTO meals (date, recipe_id)
	VALUES (?, ?)
	RETURNING meal_id;`

	var id int
	err := m.DB.QueryRow(query, date, recipeID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m MealModel) GetWeeklyByDate(date string) ([]MealsByDate, error) {
	const daysInWeek = 7

	startDate, err := getFirstDayOfWeek(date)
	if err != nil {
		return []MealsByDate{}, err
	}
	endDate := startDate.AddDate(0, 0, daysInWeek-1)

	meals := make([]MealsByDate, daysInWeek)
	mealsMap := make(map[string]*MealsByDate, daysInWeek)

	// Use slice to preserve the order of days, and Map for quick access by date
	for i := range daysInWeek {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format(DateFormat)
		meals[i] = MealsByDate{
			Date:  date,
			Meals: []Meal{},
		}
		mealsMap[dateStr] = &meals[i]
	}

	query := `SELECT m.meal_id, m.date, r.recipe_id, r.title, COALESCE(img_url, '')
	FROM meals m
	JOIN recipes r ON m.recipe_id = r.recipe_id
	WHERE m.date BETWEEN ? AND ?;`

	rows, err := m.DB.Query(query, startDate.Format(DateFormat), endDate.Format(DateFormat))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m Meal
		err := rows.Scan(&m.ID, &m.Date, &m.Recipe.ID, &m.Recipe.Title, &m.Recipe.ImageURL)
		if err != nil {
			return nil, err
		}
		mealsMap[m.Date].Meals = append(mealsMap[m.Date].Meals, m)
	}

	return meals, nil
}

func (m MealModel) RemoveByPk(id int) error {
	query := "DELETE FROM meals WHERE meal_id = ?;"

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

// getFirstDayOfWeek calculates the start of the week based on a given date.
func getFirstDayOfWeek(date string) (time.Time, error) {
	if date == "" {
		date = time.Now().Format(DateFormat)
	}

	parsed, err := time.Parse(DateFormat, date)
	if err != nil {
		return time.Time{}, err
	}

	weekday := int(parsed.Weekday())
	return parsed.AddDate(0, 0, -weekday), nil
}

package models

import (
	"database/sql"
	"time"
)

type Meal struct {
	Id      int
	Date    time.Time
	Recipes []Recipe
}

func (m Meal) DateFormat() string {
	return m.Date.Format(DateFormat)
}

func (m Meal) DayOfWeek() string {
	return m.Date.Weekday().String()
}

func (m Meal) MonthDay() string {
	return m.Date.Format("Jan 2")
}

type MealModel struct {
	DB *sql.DB
}

func (m MealModel) GetWeeklyByDate(date string) ([]Meal, error) {
	const daysInWeek = 7

	startDate, err := getFirstDayOfWeek(date)
	if err != nil {
		return []Meal{}, err
	}
	endDate := startDate.AddDate(0, 0, daysInWeek-1)

	meals := make([]Meal, daysInWeek)
	mealsMap := make(map[string]*Meal, daysInWeek)

	for i := range daysInWeek {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format(DateFormat)
		meal := Meal{
			Date:    date,
			Recipes: []Recipe{},
		}
		meals[i] = meal
		mealsMap[dateStr] = &meals[i]
	}

	query := `SELECT m.date, r.recipe_id, r.title, COALESCE(img_url, '')
	FROM meals m
	JOIN recipes r ON m.recipe_id = r.recipe_id
	WHERE m.date BETWEEN ? AND ?;`

	rows, err := m.DB.Query(query, startDate.Format(DateFormat), endDate.Format(DateFormat))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		var r Recipe

		err := rows.Scan(&date, &r.Id, &r.Title, &r.ImageURL)
		if err != nil {
			return nil, err
		}

		mealsMap[date].Recipes = append(mealsMap[date].Recipes, r)
	}

	return meals, nil
}

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

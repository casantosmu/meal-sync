package models

import (
	"database/sql"
	"errors"
	"strings"
)

var (
	ErrNotFound = errors.New("record not found")
)

const (
	DateFormat = "2006-01-02"
)

type Models struct {
	Recipe   RecipeModel
	Meal     MealModel
	Shopping ShoppingModel
}

func newNullString(s string) sql.NullString {
	if len(strings.TrimSpace(s)) == 0 {
		return sql.NullString{
			Valid: false,
		}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

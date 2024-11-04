package controllers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/casantosmu/meal-sync/models"
	"github.com/casantosmu/meal-sync/views"
)

type MealController struct {
	Logger *slog.Logger
	View   views.View
	Models models.Models
}

func (c MealController) Mount(srv *http.ServeMux) {
	srv.HandleFunc("GET /meals", c.listGET)
	srv.HandleFunc("GET /meals/recipes/selection", c.recipesSelectionGET)
}

func (c MealController) listGET(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date != "" && !isValidDate(date) {
		c.View.ClientError(w, r, http.StatusBadRequest)
		return
	}

	list, err := c.Models.Meal.GetWeeklyByDate(date)
	if err != nil {
		c.View.ServerError(w, r, err)
		return
	}

	startDate := list[0].Date

	data := map[string]any{
		"Meals":         list,
		"MonthDayYear":  startDate.Format("Jan 2, 2006"),
		"PrevWeekStart": startDate.AddDate(0, 0, -7).Format("2006-01-02"),
		"NextWeekStart": startDate.AddDate(0, 0, 7).Format("2006-01-02"),
	}
	c.View.Render(w, r, "meal-list.tmpl", data)
}

func (c MealController) recipesSelectionGET(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	search := r.URL.Query().Get("search")

	if !isValidDate(date) {
		c.View.ClientError(w, r, http.StatusBadRequest)
		return
	}

	list, err := c.Models.Recipe.Search(search)
	if err != nil {
		c.View.ServerError(w, r, err)
		return
	}

	data := map[string]any{
		"Recipes": list,
		"Date":    date,
		"Search":  search,
	}
	c.View.Partial(w, r, "meal-recipes-selection.tmpl", data)
}

func isValidDate(date string) bool {
	_, err := time.Parse(models.DateFormat, date)
	return err == nil
}

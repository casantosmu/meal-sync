package controllers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
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
	srv.HandleFunc("POST /meals", c.createPOST)
	srv.HandleFunc("GET /meals", c.listGET)
	srv.HandleFunc("DELETE /meals/{id}", c.removeDELETE)
	srv.HandleFunc("GET /meals/recipes/selection", c.recipesSelectionGET)
}

func (c MealController) createPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.View.ServerError(w, r, err)
		return
	}

	date := r.FormValue("date")
	recipeIDParam := r.FormValue("recipe_id")

	if !isValidDate(date) || recipeIDParam == "" {
		c.View.ClientError(w, r, http.StatusBadRequest)
		return
	}

	recipeID, err := strconv.Atoi(recipeIDParam)
	if err != nil {
		c.View.ClientError(w, r, http.StatusBadRequest)
		return
	}

	id, err := c.Models.Meal.Create(date, recipeID)
	if err != nil {
		c.View.ServerError(w, r, err)
		return
	}

	c.View.SetSuccessToast(w, "Added recipe to the meal plan.")

	c.Logger.Info("Meal added", "id", id)
	http.Redirect(w, r, fmt.Sprintf("/meals?date=%s", date), http.StatusSeeOther)
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
		"MealsByDate":   list,
		"MonthDayYear":  startDate.Format("Jan 2, 2006"),
		"PrevWeekStart": startDate.AddDate(0, 0, -7).Format("2006-01-02"),
		"NextWeekStart": startDate.AddDate(0, 0, 7).Format("2006-01-02"),
	}
	c.View.Render(w, r, "meal-list.tmpl", data)
}

func (c MealController) removeDELETE(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	idParam := r.PathValue("id")

	if idParam == "" {
		err := errors.New("expected id path value")
		c.View.ServerError(w, r, err)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.View.ClientError(w, r, http.StatusBadRequest)
		return
	}

	err = c.Models.Meal.RemoveByPk(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.View.ClientError(w, r, http.StatusNotFound)
			return
		}
		c.View.ServerError(w, r, err)
		return
	}

	c.View.SetSuccessToast(w, "Your recipe has been deleted.")

	c.Logger.Info("Meal deleted", "id", id)
	http.Redirect(w, r, fmt.Sprintf("/meals?date=%s", date), http.StatusSeeOther)
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

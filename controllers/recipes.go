package controllers

import (
	"log/slog"
	"net/http"

	"github.com/casantosmu/meal-sync/models"
	"github.com/casantosmu/meal-sync/views"
)

type RecipeController struct {
	Logger      *slog.Logger
	Views       views.View
	RecipeModel models.RecipeModel
}

func (c RecipeController) Mount(srv *http.ServeMux) {
	srv.HandleFunc("GET /{$}", c.listGET)
}

func (c RecipeController) listGET(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.Views.ServerError(w, r, err)
		return
	}

	search := r.FormValue("search")

	list, err := c.RecipeModel.Search(search)
	if err != nil {
		c.Views.ServerError(w, r, err)
		return
	}

	data := map[string]any{"Recipes": list, "Search": search}
	c.Views.Render(w, r, http.StatusOK, "recipe-list.tmpl", data)
}

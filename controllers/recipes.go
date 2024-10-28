package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/casantosmu/meal-sync/models"
	"github.com/casantosmu/meal-sync/views"
)

type RecipeController struct {
	Logger      *slog.Logger
	Views       views.View
	RecipeModel models.RecipeModel
}

func (c RecipeController) Mount(srv *http.ServeMux) {
	srv.HandleFunc("POST /recipes", c.createPOST)
	srv.HandleFunc("GET /{$}", c.listGET)
}

func (c RecipeController) createPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.Views.ServerError(w, r, err)
		return
	}

	title := r.FormValue("title")
	img := "https://placehold.co/300x300"

	if strings.TrimSpace(title) == "" {
		c.Views.SetErrorToast(w, "Title must not be blank.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, err := c.RecipeModel.Create(title, img)
	if err != nil {
		c.Views.ServerError(w, r, err)
		return
	}

	c.Logger.Info("Recipe created", "id", id)
	http.Redirect(w, r, fmt.Sprintf("/recipes/%d/edit", id), http.StatusSeeOther)
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

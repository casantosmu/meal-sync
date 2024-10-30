package controllers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
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
	srv.HandleFunc("GET /recipes/{id}", c.getGET)
	srv.HandleFunc("GET /recipes/{id}/edit", c.updateGET)
	srv.HandleFunc("PUT /recipes/{id}", c.updatePUT)
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

func (c RecipeController) getGET(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	if idParam == "" {
		err := errors.New("expected id path value")
		c.Views.ServerError(w, r, err)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Views.ClientError(w, r, http.StatusBadRequest)
		return
	}

	recipe, err := c.RecipeModel.GetByPk(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.Views.ClientError(w, r, http.StatusNotFound)
			return
		}
		c.Views.ServerError(w, r, err)
		return
	}

	data := map[string]any{"Recipe": recipe}
	c.Views.Render(w, r, http.StatusOK, "recipe-details.tmpl", data)
}

func (c RecipeController) updateGET(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	if idParam == "" {
		err := errors.New("expected id path value")
		c.Views.ServerError(w, r, err)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Views.ClientError(w, r, http.StatusBadRequest)
		return
	}

	recipe, err := c.RecipeModel.GetByPk(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.Views.ClientError(w, r, http.StatusNotFound)
			return
		}
		c.Views.ServerError(w, r, err)
		return
	}

	data := map[string]any{"Recipe": recipe}
	c.Views.Render(w, r, http.StatusOK, "recipe-edit.tmpl", data)
}

func (c RecipeController) updatePUT(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	if idParam == "" {
		err := errors.New("expected id path value")
		c.Views.ServerError(w, r, err)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Views.ClientError(w, r, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		c.Views.ServerError(w, r, err)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	ingredients := r.FormValue("ingredients")
	directions := r.FormValue("directions")

	validationErrs := map[string]string{}

	if strings.TrimSpace(title) == "" {
		validationErrs["title"] = "Title must not be blank."
	}

	if len(validationErrs) > 0 {
		c.Views.SetErrors(w, validationErrs)
		http.Redirect(w, r, fmt.Sprintf("/recipes/%d/edit", id), http.StatusSeeOther)
		return
	}

	err = c.RecipeModel.UpdateByPk(id, title, description, ingredients, directions)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.Views.ClientError(w, r, http.StatusNotFound)
			return
		}
		c.Views.ServerError(w, r, err)
		return
	}

	c.Views.SetSuccessToast(w, "Your recipe has been saved.")

	c.Logger.Info("Recipe updated", "id", id)
	http.Redirect(w, r, fmt.Sprintf("/recipes/%d", id), http.StatusSeeOther)
}

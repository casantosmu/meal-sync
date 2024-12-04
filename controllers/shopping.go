package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/casantosmu/meal-sync/models"
	"github.com/casantosmu/meal-sync/views"
)

type ShoppingController struct {
	Logger *slog.Logger
	View   views.View
	Models models.Models
}

// Mount registers the HTTP handlers.
func (c ShoppingController) Mount(srv *http.ServeMux) {
	srv.HandleFunc("POST /shopping/bulk", c.bulkCreatePOST)
}

func (c ShoppingController) bulkCreatePOST(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	ingredients := r.PostForm["ingredient"]

	ids, err := c.Models.Shopping.BulkCreate(ingredients)
	if err != nil {
		c.View.ServerError(w, r, err)
		return
	}

	c.View.SetSuccessToast(w, "Ingredients have been added to the sopping list.")

	c.Logger.Info("Shopping items added", "ids", ids)
	http.Redirect(w, r, fmt.Sprintf("/meals?date=%s", date), http.StatusSeeOther)
}

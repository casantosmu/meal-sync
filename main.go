package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/casantosmu/meal-sync/controllers"
	"github.com/casantosmu/meal-sync/database"
	"github.com/casantosmu/meal-sync/models"
	"github.com/casantosmu/meal-sync/views"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := database.InitDB("file:./meal_sync.db?_fk=true&_journal=WAL")
	if err != nil {
		logger.Error("Unable to connect to database", "error", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	err = database.RunMigrations(db, "./database/migrations")
	if err != nil {
		logger.Error("Migration failed", "error", err.Error())
		os.Exit(1)
	}

	view, err := views.New(logger)
	if err != nil {
		logger.Error("Unable to load templates", "error", err.Error())
		os.Exit(1)
	}

	recipeModel := models.RecipeModel{DB: db}
	recipeController := controllers.RecipeController{Logger: logger, Views: view, RecipeModel: recipeModel}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	recipeController.Mount(mux)

	logger.Info("Starting server on :3000")

	err = http.ListenAndServe(":3000", methodOverride(mux))
	if err != nil {
		logger.Error("Unable to start server", "error", err.Error())
		os.Exit(1)
	}
}

func methodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			method := r.PostFormValue("_method")
			if method == "PUT" || method == "DELETE" {
				r.Method = method
			}
		}
		next.ServeHTTP(w, r)
	})
}

package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/casantosmu/meal-sync/controllers"
	"github.com/casantosmu/meal-sync/migrations"
	"github.com/casantosmu/meal-sync/models"
	"github.com/casantosmu/meal-sync/views"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := initDB("file:./meal_sync.db?_fk=true&_journal=WAL")
	if err != nil {
		logger.Error("Unable to connect to database", "error", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	err = migrations.RunMigrations(db, "./migrations/sql")
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
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	recipeController.Mount(mux)

	logger.Info("Starting server on :3000")

	err = http.ListenAndServe(":3000", methodOverride(mux))
	if err != nil {
		logger.Error("Unable to start server", "error", err.Error())
		os.Exit(1)
	}
}

func initDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
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

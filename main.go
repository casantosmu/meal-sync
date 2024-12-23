package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"slices"

	"github.com/casantosmu/meal-sync/controllers"
	"github.com/casantosmu/meal-sync/middlewares"
	"github.com/casantosmu/meal-sync/migrations"
	"github.com/casantosmu/meal-sync/models"
	"github.com/casantosmu/meal-sync/views"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := initDB("file:./meal_sync.db?_fk=true&_journal=WAL") // Enable foreign keys and use WAL for better concurrency
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

	models := models.Models{
		Recipe:   models.RecipeModel{DB: db},
		Meal:     models.MealModel{DB: db},
		Shopping: models.ShoppingModel{DB: db},
	}

	srv := buildServer(
		[]controller{
			controllers.RecipeController{Logger: logger, View: view, Models: models},
			controllers.MealController{Logger: logger, View: view, Models: models},
			controllers.ShoppingController{Logger: logger, View: view, Models: models},
		},
		[]middleware{
			middlewares.RecoverPanic(view),
			middlewares.MethodOverride,
			middlewares.LogRequest(logger),
			middlewares.Security,
		})

	logger.Info("Starting server on :3000")

	err = http.ListenAndServe(":3000", srv)
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

type controller interface {
	Mount(mux *http.ServeMux)
}

type middleware func(next http.Handler) http.Handler

func buildServer(controllers []controller, middlewares []middleware) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	for _, controller := range controllers {
		controller.Mount(mux)
	}

	slices.Reverse(middlewares) // Reverse the middleware slice to maintain the intended order when chaining
	var next http.Handler = mux
	for _, middleware := range middlewares {
		next = middleware(next)
	}

	return next
}

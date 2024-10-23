package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/casantosmu/meal-sync/database"
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

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Message string
		}{
			Message: "Hello world!!",
		}

		tmpl := template.Must(template.ParseFiles("./views/index.tmpl"))
		err := tmpl.Execute(w, data)
		if err != nil {
			logger.Error(err.Error())
		}
	})

	logger.Info("Starting server on :3000")

	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		logger.Error("Unable to start server", "error", err.Error())
		os.Exit(1)
	}
}

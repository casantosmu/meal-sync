package views

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"text/template"
)

type View struct {
	logger    *slog.Logger
	templates map[string]*template.Template
}

func New(logger *slog.Logger) (View, error) {
	pages, err := filepath.Glob("./views/pages/*.tmpl")
	if err != nil {
		return View{}, err
	}

	templates := make(map[string]*template.Template, len(pages))
	for _, page := range pages {
		name := filepath.Base(page)
		tmpl, err := template.New(name).ParseFiles("./views/layouts/base.tmpl", page)
		if err != nil {
			return View{}, err
		}
		templates[name] = tmpl
	}

	return View{logger: logger, templates: templates}, nil
}

func (v View) Render(w http.ResponseWriter, r *http.Request, status int, page string, data map[string]any) {
	flash, err := getFlash(w, r)
	if err != nil {
		v.ServerError(w, r, err)
		return
	}
	data["Toast"] = flash.Toast
	data["Errors"] = flash.Errors

	tmpl, ok := v.templates[page]
	if !ok {
		err := fmt.Errorf("template '%s' does not exist", page)
		v.ServerError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		v.ServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = buf.WriteTo(w); err != nil {
		v.logger.Warn("Error writing response", "error", err.Error())
	}
}

func (t View) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	t.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (t View) ClientError(w http.ResponseWriter, r *http.Request, status int) {
	http.Error(w, http.StatusText(status), status)
}

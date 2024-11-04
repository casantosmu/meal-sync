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
	logger   *slog.Logger
	pages    map[string]*template.Template
	partials map[string]*template.Template
}

func New(logger *slog.Logger) (View, error) {
	pagesFiles, err := filepath.Glob("./views/pages/*.tmpl")
	if err != nil {
		return View{}, err
	}
	partialsFiles, err := filepath.Glob("./views/partials/*.tmpl")
	if err != nil {
		return View{}, err
	}

	pages := make(map[string]*template.Template, len(pagesFiles))
	for _, page := range pagesFiles {
		name := filepath.Base(page)
		tmpl, err := template.New(name).ParseFiles("./views/layouts/base.tmpl", page)
		if err != nil {
			return View{}, err
		}
		pages[name] = tmpl
	}

	partials := make(map[string]*template.Template, len(partialsFiles))
	for _, partial := range partialsFiles {
		name := filepath.Base(partial)
		tmpl, err := template.New(name).ParseFiles(partial)
		if err != nil {
			return View{}, err
		}
		partials[name] = tmpl
	}

	return View{logger: logger, pages: pages, partials: partials}, nil
}

func (v View) Render(w http.ResponseWriter, r *http.Request, page string, data map[string]any) {
	flash, err := getFlash(w, r)
	if err != nil {
		v.ServerError(w, r, err)
		return
	}
	data["Toast"] = flash.Toast
	data["Errors"] = flash.Errors

	tmpl, ok := v.pages[page]
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

func (v View) Partial(w http.ResponseWriter, r *http.Request, partial string, data map[string]any) {
	tmpl, ok := v.partials[partial]
	if !ok {
		err := fmt.Errorf("partial template '%s' does not exist", partial)
		v.ServerError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, data)
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

package views

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"text/template"
)

type contextKey string

const NonceKey contextKey = "nonce"

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

		tmpl, err := template.New(name).ParseFiles("./views/layouts/base.tmpl")
		if err != nil {
			return View{}, err
		}

		tmpl, err = tmpl.ParseGlob("./views/templates/*.tmpl")
		if err != nil {
			return View{}, err
		}

		tmpl, err = tmpl.ParseFiles(page)
		if err != nil {
			return View{}, err
		}

		pages[name] = tmpl
	}

	partials := make(map[string]*template.Template, len(partialsFiles))
	for _, partial := range partialsFiles {
		name := filepath.Base(partial)

		tmpl, err := template.New(name).ParseFiles("./views/layouts/partial.tmpl")
		if err != nil {
			return View{}, err
		}

		tmpl, err = tmpl.ParseGlob("./views/templates/*.tmpl")
		if err != nil {
			return View{}, err
		}

		tmpl, err = tmpl.ParseFiles(partial)
		if err != nil {
			return View{}, err
		}

		partials[name] = tmpl
	}

	return View{
		logger:   logger,
		pages:    pages,
		partials: partials,
	}, nil
}

// Render generates a complete HTML page with layout and flash messages.
// Used for full-page responses requiring a complete HTML document.
func (v View) Render(w http.ResponseWriter, r *http.Request, page string, data map[string]any) {
	flash, err := getFlash(w, r)
	if err != nil {
		v.ServerError(w, r, err)
		return
	}
	data["Toast"] = flash.Toast
	data["Errors"] = flash.Errors

	nonce, ok := r.Context().Value(NonceKey).(string)
	if ok {
		data["Nonce"] = nonce
	}

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

// Partial renders an HTML fragment (partial) without a full layout and includes flash messages.
// Used for sections of a page, often in AJAX responses.
func (v View) Partial(w http.ResponseWriter, r *http.Request, partial string, data map[string]any) {
	flash, err := getFlash(w, r)
	if err != nil {
		v.ServerError(w, r, err)
		return
	}
	data["Toast"] = flash.Toast
	data["Errors"] = flash.Errors

	nonce, ok := r.Context().Value(NonceKey).(string)
	if ok {
		data["Nonce"] = nonce
	}

	tmpl, ok := v.partials[partial]
	if !ok {
		err := fmt.Errorf("partial template '%s' does not exist", partial)
		v.ServerError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, "partial", data)
	if err != nil {
		v.ServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = buf.WriteTo(w); err != nil {
		v.logger.Warn("Error writing response", "error", err.Error())
	}
}

// ServerError sends a 500 error response and logs the error.
func (t View) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	t.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// ClientError sends a specified HTTP error status to the user.
func (t View) ClientError(w http.ResponseWriter, r *http.Request, status int) {
	http.Error(w, http.StatusText(status), status)
}

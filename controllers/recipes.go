package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/casantosmu/meal-sync/models"
	"github.com/casantosmu/meal-sync/views"
)

const maxFileSize = 1 * 1024 * 1024 // 1 MB

var ErrFileFormat = errors.New("Unsupported file format.")

type RecipeController struct {
	Logger *slog.Logger
	View   views.View
	Models models.Models
}

// Mount registers the HTTP handlers.
func (c RecipeController) Mount(srv *http.ServeMux) {
	srv.HandleFunc("POST /recipes", c.createPOST)
	srv.HandleFunc("GET /{$}", c.listGET)
	srv.HandleFunc("GET /recipes/{id}", c.getGET)
	srv.HandleFunc("GET /recipes/{id}/edit", c.updateGET)
	srv.HandleFunc("PUT /recipes/{id}", c.updatePUT)
	srv.HandleFunc("DELETE /recipes/{id}", c.removeDELETE)

	srv.HandleFunc("GET /recipes/{id}/image", c.imageGET)
	srv.HandleFunc("PUT /recipes/{id}/image", c.imagePUT)
	srv.HandleFunc("DELETE /recipes/{id}/image", c.imageDELETE)
}

func (c RecipeController) createPOST(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")

	if strings.TrimSpace(title) == "" {
		c.View.SetErrorToast(w, "Title must not be blank.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, err := c.Models.Recipe.Create(title)
	if err != nil {
		c.View.ServerError(w, r, err)
		return
	}

	c.Logger.Info("Recipe created", "id", id)
	http.Redirect(w, r, fmt.Sprintf("/recipes/%d/edit", id), http.StatusSeeOther)
}

func (c RecipeController) listGET(w http.ResponseWriter, r *http.Request) {
	search := r.PostFormValue("search")

	list, err := c.Models.Recipe.Search(search)
	if err != nil {
		c.View.ServerError(w, r, err)
		return
	}

	data := map[string]any{"Recipes": list, "Search": search}
	c.View.Render(w, r, "recipe-list.tmpl", data)
}

func (c RecipeController) getGET(w http.ResponseWriter, r *http.Request) {
	id, err := pathInt(r, "id")
	if err != nil {
		if errors.Is(err, ErrPathValueNotFound) {
			c.View.ServerError(w, r, err)
			return
		}
		c.View.ClientError(w, r, http.StatusBadRequest)
		return
	}

	recipe, err := c.Models.Recipe.GetByPk(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.View.ClientError(w, r, http.StatusNotFound)
			return
		}
		c.View.ServerError(w, r, err)
		return
	}

	data := map[string]any{"Recipe": recipe}
	c.View.Render(w, r, "recipe-details.tmpl", data)
}

func (c RecipeController) updateGET(w http.ResponseWriter, r *http.Request) {
	id, err := pathInt(r, "id")
	if err != nil {
		if errors.Is(err, ErrPathValueNotFound) {
			c.View.ServerError(w, r, err)
			return
		}
		c.View.ClientError(w, r, http.StatusBadRequest)
		return
	}

	recipe, err := c.Models.Recipe.GetByPk(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.View.ClientError(w, r, http.StatusNotFound)
			return
		}
		c.View.ServerError(w, r, err)
		return
	}

	data := map[string]any{"Recipe": recipe}
	c.View.Render(w, r, "recipe-edit.tmpl", data)
}

func (c RecipeController) updatePUT(w http.ResponseWriter, r *http.Request) {
	id, err := pathInt(r, "id")
	if err != nil {
		if errors.Is(err, ErrPathValueNotFound) {
			c.View.ServerError(w, r, err)
			return
		}
		c.View.ClientError(w, r, http.StatusBadRequest)
		return
	}

	title := r.PostFormValue("title")
	description := r.PostFormValue("description")
	ingredients := r.PostFormValue("ingredients")
	directions := r.PostFormValue("directions")

	validationErrs := map[string]string{}

	if strings.TrimSpace(title) == "" {
		validationErrs["title"] = "Title must not be blank."
	}

	if len(validationErrs) > 0 {
		c.View.SetErrors(w, validationErrs)
		http.Redirect(w, r, fmt.Sprintf("/recipes/%d/edit", id), http.StatusSeeOther)
		return
	}

	err = c.Models.Recipe.UpdateByPk(id, title, description, ingredients, directions)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.View.ClientError(w, r, http.StatusNotFound)
			return
		}
		c.View.ServerError(w, r, err)
		return
	}

	c.View.SetSuccessToast(w, "Your recipe has been saved.")

	c.Logger.Info("Recipe updated", "id", id)
	http.Redirect(w, r, fmt.Sprintf("/recipes/%d", id), http.StatusSeeOther)
}

func (c RecipeController) removeDELETE(w http.ResponseWriter, r *http.Request) {
	id, err := pathInt(r, "id")
	if err != nil {
		if errors.Is(err, ErrPathValueNotFound) {
			c.View.ServerError(w, r, err)
			return
		}
		c.View.ClientError(w, r, http.StatusBadRequest)
		return
	}

	err = c.Models.Recipe.RemoveByPk(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.View.ClientError(w, r, http.StatusNotFound)
			return
		}
		c.View.ServerError(w, r, err)
		return
	}

	c.View.SetSuccessToast(w, "Your recipe has been deleted.")

	c.Logger.Info("Recipe deleted", "id", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c RecipeController) imageGET(w http.ResponseWriter, r *http.Request) {
	id, err := pathInt(r, "id")
	if err != nil {
		if errors.Is(err, ErrPathValueNotFound) {
			c.View.ServerError(w, r, err)
			return
		}
		c.View.ClientError(w, r, http.StatusBadRequest)
		return
	}

	recipe, err := c.Models.Recipe.GetByPk(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.View.ClientError(w, r, http.StatusNotFound)
			return
		}
		c.View.ServerError(w, r, err)
		return
	}

	data := map[string]any{"Recipe": recipe}
	c.View.Partial(w, r, "recipe-image-upload.tmpl", data)
}

func (c RecipeController) imagePUT(w http.ResponseWriter, r *http.Request) {
	id, err := pathInt(r, "id")
	if err != nil {
		if errors.Is(err, ErrPathValueNotFound) {
			c.View.ServerError(w, r, err)
			return
		}
		c.View.ClientError(w, r, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			c.View.SetErrorToast(w, "Please upload an image file.")
			http.Redirect(w, r, fmt.Sprintf("/recipes/%d/image", id), http.StatusSeeOther)
			return
		}
		c.View.ServerError(w, r, err)
		return
	}
	defer file.Close()

	if handler.Size > maxFileSize {
		c.View.SetErrorToast(w, "File size exceeds 1 MB. Please upload a smaller file.")
		http.Redirect(w, r, fmt.Sprintf("/recipes/%d/image", id), http.StatusSeeOther)
		return
	}

	ext := filepath.Ext(handler.Filename)

	err = validateFileFormat(file, ext)
	if err != nil {
		if errors.Is(err, ErrFileFormat) {
			c.View.SetErrorToast(w, "Unsupported file format. Please upload a .jpg, .jpeg, or .png file.")
			http.Redirect(w, r, fmt.Sprintf("/recipes/%d/image", id), http.StatusSeeOther)
			return
		}
		c.View.ServerError(w, r, err)
		return
	}

	path, err := uploadImage(file, ext)
	if err != nil {
		c.View.ServerError(w, r, err)
		return
	}

	err = c.Models.Recipe.UpdateImageByPk(id, path)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.View.ClientError(w, r, http.StatusNotFound)
			return
		}
		c.View.ServerError(w, r, err)
		return
	}

	c.View.SetSuccessToast(w, "Your image has been uploaded.")

	c.Logger.Info("Image updated", "id", id, "path", path)
	http.Redirect(w, r, fmt.Sprintf("/recipes/%d/image", id), http.StatusSeeOther)
}

func (c RecipeController) imageDELETE(w http.ResponseWriter, r *http.Request) {
	id, err := pathInt(r, "id")
	if err != nil {
		if errors.Is(err, ErrPathValueNotFound) {
			c.View.ServerError(w, r, err)
			return
		}
		c.View.ClientError(w, r, http.StatusBadRequest)
		return
	}

	recipe, err := c.Models.Recipe.GetByPk(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.View.ClientError(w, r, http.StatusNotFound)
			return
		}
		c.View.ServerError(w, r, err)
		return
	}

	err = c.Models.Recipe.UpdateImageByPk(id, "")
	if err != nil {
		c.View.ServerError(w, r, err)
		return
	}

	path := strings.TrimPrefix(recipe.ImageURL, "/")
	if path != "" {
		err := os.Remove(path)
		if err != nil {
			c.View.ServerError(w, r, err)
			return
		}
	}

	c.View.SetSuccessToast(w, "The image has been deleted.")

	c.Logger.Info("Image deleted", "id", id, "path", path)
	http.Redirect(w, r, fmt.Sprintf("/recipes/%d/image", id), http.StatusSeeOther)
}

func validateFileFormat(file multipart.File, ext string) error {
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	if !validExtensions[ext] {
		return ErrFileFormat
	}

	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	mimeType := http.DetectContentType(buf)
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		return ErrFileFormat
	}

	return nil
}

func generateFilename(ext string) (string, error) {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf) + ext, nil
}

func uploadImage(file multipart.File, ext string) (string, error) {
	filename, err := generateFilename(ext)
	if err != nil {
		return "", err
	}

	path := filepath.Join("./uploads/images", filename)

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return "/" + path, nil
}

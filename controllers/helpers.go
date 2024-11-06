package controllers

import (
	"errors"
	"net/http"
	"strconv"
)

var (
	ErrPathValueNotFound  = errors.New("expected path value is missing")
	ErrPathValueNotNumber = errors.New("path value is not a valid number")
)

// pathInt parses and validates an integer parameter from the request path.
// Returns an error if the value is missing or not a valid integer.
func pathInt(r *http.Request, name string) (int, error) {
	value := r.PathValue(name)
	if value == "" {
		return 0, ErrPathValueNotFound
	}

	int, err := strconv.Atoi(value)
	if err != nil {
		return 0, ErrPathValueNotNumber
	}

	return int, nil
}

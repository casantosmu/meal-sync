package views

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

const (
	successName = "flash_success"
	errorName   = "flash_error"
	errorsName  = "flash_errors"
)

type toast struct {
	Success string
	Error   string
}

type flash struct {
	Toast  toast
	Errors map[string]string
}

func (v View) SetSuccessToast(w http.ResponseWriter, value string) {
	setFlashValue(w, successName, []byte(value))
}

func (v View) SetErrorToast(w http.ResponseWriter, value string) {
	setFlashValue(w, errorName, []byte(value))
}

func (v View) SetErrors(w http.ResponseWriter, errs map[string]string) {
	serialized, err := json.Marshal(errs)
	if err != nil {
		v.logger.Warn("Failed to serialize flash errors", "error", err.Error())
		return
	}
	setFlashValue(w, errorsName, serialized)
}

func (v View) SetFlashBool(w http.ResponseWriter, name string, value bool) {
	str := "true"
	if !value {
		str = "false"
	}
	setFlashValue(w, name, []byte(str))
}

func (v View) GetFlashBool(w http.ResponseWriter, r *http.Request, name string) (bool, error) {
	value, err := getFlashValue(w, r, name)
	if err != nil {
		return false, err
	}
	return string(value) == "true", nil
}

func setFlashValue(w http.ResponseWriter, name string, value []byte) {
	encoded := base64.URLEncoding.EncodeToString(value)
	c := &http.Cookie{Name: name, Value: encoded, Path: "/", MaxAge: 15}
	http.SetCookie(w, c)
}

func getFlashValue(w http.ResponseWriter, r *http.Request, name string) ([]byte, error) {
	c, err := r.Cookie(name)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return nil, nil
		default:
			return nil, err
		}
	}

	decoded, err := base64.URLEncoding.DecodeString(c.Value)
	if err != nil {
		return nil, err
	}

	dc := &http.Cookie{Name: name, Path: "/", MaxAge: -1}
	http.SetCookie(w, dc)

	return decoded, nil
}

func getFlash(w http.ResponseWriter, r *http.Request) (flash, error) {
	successDecoded, err := getFlashValue(w, r, successName)
	if err != nil {
		return flash{}, err
	}

	errorDecoded, err := getFlashValue(w, r, errorName)
	if err != nil {
		return flash{}, err
	}

	toast := toast{
		Success: string(successDecoded),
		Error:   string(errorDecoded),
	}

	errorsDecoded, err := getFlashValue(w, r, errorsName)
	if err != nil {
		return flash{}, err
	}

	// If errors is nil do not unserialize
	if errorsDecoded == nil {
		return flash{Toast: toast}, nil
	}

	var errorsMap map[string]string
	err = json.Unmarshal(errorsDecoded, &errorsMap)
	if err != nil {
		return flash{}, err
	}

	return flash{Toast: toast, Errors: errorsMap}, nil
}

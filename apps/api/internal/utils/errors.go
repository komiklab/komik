package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
)

type KomikError struct {
	ErrorCode       string         `json:"error_code"`
	ErrorMessage    string         `json:"error_message"`
	DetailedMessage string         `json:"-"`
	OriginalErrors  error          `json:"-"`
	StatusCode      http.ConnState `json:"-"`
}

func NewKomikError() *KomikError {
	return &KomikError{}
}

func (e *KomikError) Error() string {
	log.Error().Str("error_code", e.ErrorCode).Msg(e.ErrorMessage)
	jsonData, err := json.Marshal(e)
	if err != nil {
		return `{"error_code":"unknown_error","error_message":"could not marshall error message"}`
	}
	return string(jsonData)
}

func (e *KomikError) AddOriginalError(err error) {
	e.OriginalErrors = errors.Join(err, e.OriginalErrors)
	e.DetailedMessage = e.DetailedMessage + ":: " + err.Error()
}

func (e *KomikError) clone() *KomikError {
	return &KomikError{
		ErrorCode:       e.ErrorCode,
		ErrorMessage:    e.ErrorMessage,
		DetailedMessage: e.DetailedMessage,
		OriginalErrors:  e.OriginalErrors,
		StatusCode:      e.StatusCode,
	}
}

func (e *KomikError) WithDetailedMessage(message string) *KomikError {
	clone := e.clone()
	clone.DetailedMessage = message
	return clone
}

func (e *KomikError) WithOriginalError(err error) *KomikError {
	clone := e.clone()
	clone.AddOriginalError(err)
	return clone
}

func (e *KomikError) WithStatusCode(statusCode http.ConnState) *KomikError {
	clone := e.clone()
	clone.StatusCode = statusCode
	return clone
}

func (e *KomikError) WithErrorCode(code string) *KomikError {
	clone := e.clone()
	clone.ErrorCode = code
	return clone
}

func (e *KomikError) WithErrorMessage(message string) *KomikError {
	clone := e.clone()
	clone.ErrorMessage = message
	return clone
}

func (e *KomikError) Wrap(err error) *KomikError {
	clone := e.clone()
	clone.AddOriginalError(err)
	return clone
}

// error factory functions

func NewBindError(message string, err error) *KomikError {
	return NewKomikError().
		WithErrorCode("bind_error").
		WithErrorMessage(message).
		WithOriginalError(err)
}

func NewDatabaseError(message string, err error) *KomikError {
	return NewKomikError().
		WithErrorCode("database_error").
		WithErrorMessage(message).
		WithStatusCode(http.StatusInternalServerError).
		WithOriginalError(err)
}

func NewValidationError(message string, err error) *KomikError {
	return NewKomikError().
		WithErrorCode("validation_error").
		WithErrorMessage(message).
		WithStatusCode(http.StatusBadRequest).
		WithOriginalError(err)
}

func NewAuthenticationError(message string, err error) *KomikError {
	return NewKomikError().
		WithErrorCode("authentication_error").
		WithErrorMessage(message).
		WithStatusCode(http.StatusUnauthorized).
		WithOriginalError(err)
}

func NewAuthorizationError(message string, err error) *KomikError {
	return NewKomikError().
		WithErrorCode("authorization_error").
		WithErrorMessage(message).
		WithStatusCode(http.StatusForbidden).
		WithOriginalError(err)
}

func NewNotFoundError(message string, err error) *KomikError {
	return NewKomikError().
		WithErrorCode("not_found_error").
		WithErrorMessage(message).
		WithStatusCode(http.StatusNotFound).
		WithOriginalError(err)
}

func NewConflictError(message string, err error) *KomikError {
	return NewKomikError().
		WithErrorCode("conflict_error").
		WithErrorMessage(message).
		WithStatusCode(http.StatusConflict).
		WithOriginalError(err)
}

func NewInternalServerError(message string, err error) *KomikError {
	return NewKomikError().
		WithErrorCode("internal_server_error").
		WithErrorMessage(message).
		WithStatusCode(http.StatusInternalServerError).
		WithOriginalError(err)
}

func NewBadRequestError(message string, err error) *KomikError {
	return NewKomikError().
		WithErrorCode("bad_request_error").
		WithErrorMessage(message).
		WithStatusCode(http.StatusBadRequest).
		WithOriginalError(err)
}

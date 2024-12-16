package pkgtypes

import (
	"errors"
	"fmt"
)

// ErrorType define los tipos de errores base
type ErrorType string

// Constantes para ErrorType
const (
	ErrNotFound        ErrorType = "NOT_FOUND"
	ErrConflict        ErrorType = "CONFLICT"
	ErrInvalidInput    ErrorType = "INVALID_INPUT"
	ErrOperationFailed ErrorType = "OPERATION_FAILED"
	ErrValidation      ErrorType = "VALIDATION_ERROR"
	ErrConnection      ErrorType = "CONNECTION_ERROR"
	ErrTimeout         ErrorType = "TIMEOUT"
	ErrUnavailable     ErrorType = "SERVICE_UNAVAILABLE"
	ErrAuthentication  ErrorType = "AUTHENTICATION_ERROR"
	ErrAuthorization   ErrorType = "AUTHORIZATION_ERROR"
)

// Error representa un error del dominio
type Error struct {
	Type    ErrorType      `json:"type"`
	Message string         `json:"message"`
	Details error          `json:"-"`
	Context map[string]any `json:"context,omitempty"`
}

// Métodos para Error
func (e *Error) Error() string {
	if e.Details != nil {
		return fmt.Sprintf("%s: %s (details: %v)", e.Type, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Details
}

// Constructores para Error
func NewError(errType ErrorType, message string, details error) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Details: details,
	}
}

func NewErrorWithContext(errType ErrorType, message string, details error, context map[string]any) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Details: details,
		Context: context,
	}
}

// Funciones helper para verificar tipos de error
func IsNotFound(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Type == ErrNotFound
}

func IsConflict(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Type == ErrConflict
}

func IsValidationError(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Type == ErrValidation
}

// GetErrorType extrae el tipo de error
func GetErrorType(err error) (ErrorType, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e.Type, true
	}
	return "", false
}

// GetErrorContext obtiene el contexto del error
func GetErrorContext(err error) (map[string]any, bool) {
	var e *Error
	if errors.As(err, &e) && e.Context != nil {
		return e.Context, true
	}
	return nil, false
}

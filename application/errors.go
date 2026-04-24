package application

import "fmt"

// ValidationError representa un error de validación
type ValidationError struct {
	message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.message)
}

func NewValidationError(message string) ValidationError {
	return ValidationError{message: message}
}

// DomainError representa un error de negocio
type DomainError struct {
	message string
}

func (e DomainError) Error() string {
	return fmt.Sprintf("domain error: %s", e.message)
}

func NewDomainError(message string) DomainError {
	return DomainError{message: message}
}

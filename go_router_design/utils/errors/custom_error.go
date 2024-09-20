package errors

// NewCustomError creates a new custom error
func NewCustomError(status int, message string) *CustomError {
	return &CustomError{
		code:    status,
		message: message,
	}
}

// CustomError es una interfaz para definir errores custom
type CustomError struct {
	code    int
	message string
}

// Code http error code
func (e *CustomError) Code() int {
	return e.code
}

// Message http error message
func (e *CustomError) Error() string {
	return e.message
}

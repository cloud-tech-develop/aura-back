package errors

import "errors"

var (
	ErrUnauthorized          = errors.New("unauthorized")
	ErrInvalidCredentials    = errors.New("credenciales inválidas")
	ErrEnterpriseNotFound    = errors.New("enterprise no encontrado")
	ErrEnterpriseInactive    = errors.New("enterprise inactiva")
	ErrPasswordRequired      = errors.New("password requerido")
	ErrFieldsRequired        = errors.New("name, slug y email requeridos")
	ErrInvalidSlug           = errors.New("slug inválido: solo minúsculas, números y _")
	ErrTokenGeneration       = errors.New("error al generar token")
	ErrEventPublish          = errors.New("error al publicar evento")
	ErrEventNil              = errors.New("event cannot be nil")
	ErrEventNameEmpty        = errors.New("event name cannot be empty")
	ErrHandlerNil            = errors.New("handler cannot be nil")
	ErrTimeoutPublishing     = errors.New("timeout publishing event")
	ErrNoHandlers            = errors.New("no handlers for event")
	ErrInvalidPayload        = errors.New("invalid payload type")
	ErrOpenLogFile           = errors.New("failed to open log file")
	ErrMarshalLogEntry       = errors.New("failed to marshal log entry")
	ErrAuthHeaderRequired    = errors.New("authorization header requerido")
	ErrTokenInvalid          = errors.New("token inválido")
	ErrIPNotValid            = errors.New("ip no válida")
	ErrEmailPasswordRequired = errors.New("email y password requeridos")
)

type ErrorCode string

const (
	CodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	CodeInvalidCreds  ErrorCode = "INVALID_CREDENTIALS"
	CodeNotFound      ErrorCode = "NOT_FOUND"
	CodeForbidden     ErrorCode = "FORBIDDEN"
	CodeBadRequest    ErrorCode = "BAD_REQUEST"
	CodeInternalError ErrorCode = "INTERNAL_ERROR"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

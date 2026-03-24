package errors

import "errors"

var (
	// Auth
	ErrAuthHeaderRequired    = errors.New("Authorization header required")
	ErrTokenInvalid          = errors.New("Token inválido o expirado")
	ErrTokenGeneration       = errors.New("Error generando token")
	ErrIPNotValid            = errors.New("IP no válida para este token")
	ErrEmailPasswordRequired = errors.New("Email y password son requeridos")
	ErrInvalidCredentials    = errors.New("Credenciales inválidas")

	// Enterprise
	ErrEnterpriseInactive = errors.New("Empresa inactiva")
	ErrEnterpriseNotFound = errors.New("Empresa no encontrada")
	ErrUserInactive       = errors.New("Usuario inactivo")

	// Generic
	ErrFieldsRequired = errors.New("Campos requeridos faltantes o inválidos")
	ErrNotFound       = errors.New("Recurso no encontrado")
	ErrInternal       = errors.New("Error interno del servidor")
	ErrForbidden      = errors.New("Acceso denegado")

	// VO
	ErrEmailRequired    = "Email es requerido"
	ErrInvalidEmail     = "Formato de email inválido"
	ErrDocumentRequired = "Documento es requerido"
	ErrInvalidDocument  = "Documento debe contener entre 5 y 20 caracteres alfanuméricos o guiones"

	InternalError = "Error interno del servidor"
)

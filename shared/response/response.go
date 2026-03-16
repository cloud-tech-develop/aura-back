package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	_CREATED = "Creado exitosamente"
	_UPDATED = "Actualizado exitosamente"
	_DELETED = "Eliminado exitosamente"
	_SUCCESS = "Operación exitosa"
)

// OK sends a 200 JSON response.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data, "success": true, "message": _SUCCESS})
}

// Created sends a 201 JSON response.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{"data": data, "success": true, "message": _CREATED})
}

// Updated sends a 200 JSON response.
func Updated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data, "success": true, "message": _UPDATED})
}

// Deleted sends a 200 JSON response.
func Deleted(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": nil, "success": true, "message": _DELETED})
}

// BadRequest sends a 400 JSON error response.
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"message": msg, "success": false, "data": nil})
}

// Unauthorized sends a 401 JSON error response.
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, gin.H{"message": msg, "success": false, "data": nil})
}

// Forbidden sends a 403 JSON error response.
func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, gin.H{"message": msg, "success": false, "data": nil})
}

// NotFound sends a 404 JSON error response.
func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, gin.H{"message": msg, "success": false, "data": nil})
}

// Conflict sends a 409 JSON error response.
func Conflict(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, gin.H{"error": "conflict", "message": msg, "success": false, "data": nil})
}

// Internal sends a 500 JSON error response.
func InternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, gin.H{"message": msg, "success": false, "data": nil})
}

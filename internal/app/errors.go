package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHTML aborts the request and sends an HTML error response.
// Used by v0 endpoints that return HTML for HTMX.
func ErrorHTML(ctx *gin.Context, status int, msg string) {
	ctx.Data(status, "text/html; charset=utf-8", []byte(msg))
	ctx.Abort()
}

// ErrorJSON aborts the request and sends a JSON error response.
// Used by v1 endpoints that return JSON.
func ErrorJSON(ctx *gin.Context, status int, msg string) {
	ctx.AbortWithStatusJSON(status, gin.H{
		"message": msg,
		"status":  status,
	})
}

// ErrorText aborts the request and sends a plain text error response.
// Used by login/auth endpoints.
func ErrorText(ctx *gin.Context, status int, msg string) {
	ctx.String(status, msg)
	ctx.Abort()
}

// ErrorInternalHTML is a shorthand for a 500 HTML error.
func ErrorInternalHTML(ctx *gin.Context, msg string) {
	ErrorHTML(ctx, http.StatusInternalServerError, msg)
}

// ErrorBadRequestHTML is a shorthand for a 400 HTML error.
func ErrorBadRequestHTML(ctx *gin.Context, msg string) {
	ErrorHTML(ctx, http.StatusBadRequest, msg)
}

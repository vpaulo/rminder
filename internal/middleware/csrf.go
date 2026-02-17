package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const csrfTokenKey = "csrf_token"
const csrfHeaderName = "X-CSRF-Token"

func generateCSRFToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func CSRFMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		// Ensure a CSRF token exists in the session
		token, _ := session.Get(csrfTokenKey).(string)
		if token == "" {
			var err error
			token, err = generateCSRFToken()
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			session.Set(csrfTokenKey, token)
			if err := session.Save(); err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		// Store token in context for templates
		ctx.Set(csrfTokenKey, token)

		// Validate token on state-changing methods
		if ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" || ctx.Request.Method == "DELETE" {
			clientToken := ctx.GetHeader(csrfHeaderName)
			if clientToken == "" || clientToken != token {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		ctx.Next()
	}
}

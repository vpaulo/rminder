package middleware

import "github.com/gin-gonic/gin"

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip security headers for auth routes that redirect to external providers
		path := ctx.Request.URL.Path
		if path == "/login" || path == "/callback" || path == "/logout" {
			ctx.Next()
			return
		}

		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-XSS-Protection", "0")
		ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		ctx.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-eval'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; img-src 'self' data:; font-src 'self' https://fonts.gstatic.com; connect-src 'self'; frame-ancestors 'none'")
		ctx.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		ctx.Next()
	}
}

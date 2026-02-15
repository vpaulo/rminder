package app

import "github.com/gin-gonic/gin"

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-XSS-Protection", "0")
		ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		ctx.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'")
		ctx.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		ctx.Next()
	}
}

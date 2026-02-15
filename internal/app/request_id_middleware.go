package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey = "request_id"

func GetRequestID(ctx *gin.Context) string {
	return ctx.GetString(requestIDKey)
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Set(requestIDKey, requestID)
		ctx.Header("X-Request-ID", requestID)

		ctx.Next()
	}
}

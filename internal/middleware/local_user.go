package middleware

import (
	"rminder/internal/app"
	"rminder/internal/app/database"

	"github.com/gin-gonic/gin"
)

// LocalMiddleware injects a database directly from a path, bypassing user/session auth.
func LocalMiddleware(dbPath string, s *app.App) gin.HandlerFunc {
	db := database.New(dbPath)
	return func(ctx *gin.Context) {
		log := s.Logger().WithRequestID(app.GetRequestID(ctx))
		app.SetUserDatabase(ctx, db)
		app.SetLogger(ctx, log)
		ctx.Next()
	}
}

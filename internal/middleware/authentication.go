package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"rminder/internal/app"
	"rminder/internal/database"
	"rminder/internal/user"
)

func SetUserDatabase(ctx *gin.Context, db database.Service) {
	ctx.Set("user_database", db)
}

func GetUserDatabase(ctx *gin.Context) database.Service {
	return ctx.MustGet("user_database").(database.Service)
}

func Authentication(s *app.App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		user_id := user.GetUserId(session)
		if user_id == "" {
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}

		db, err := s.GetDatabaseForUser(user_id)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		SetUserDatabase(ctx, db)

		ctx.Next()

	}
}

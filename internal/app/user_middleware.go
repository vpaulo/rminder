package app

import (
	"net/http"

	"rminder/internal/app/database"
	"rminder/internal/app/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetUserDatabase(ctx *gin.Context, db database.Service) {
	ctx.Set("user_database", db)
}

func GetUserDatabase(ctx *gin.Context) database.Service {
	return ctx.MustGet("user_database").(database.Service)
}

func SetUser(ctx *gin.Context, user *user.User) {
	ctx.Set("user", user)
}

func GetUser(ctx *gin.Context) *user.User {
	return ctx.MustGet("user").(*user.User)
}

func UserMiddleware(s *App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		user_id := user.GetUserId(session)
		if user_id == "" {
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}

		user, err := s.GetUser(user_id)
		if err != nil {
			ctx.Redirect(http.StatusSeeOther, "/logout")
			return
		}
		SetUser(ctx, user)

		db, err := s.GetDatabaseForUser(user_id)
		if err != nil {
			ctx.Redirect(http.StatusSeeOther, "/logout")
			return
		}
		SetUserDatabase(ctx, db)

		ctx.Next()

	}
}

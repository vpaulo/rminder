package app

import (
	"net/http"

	"rminder/internal/app/database"
	"rminder/internal/app/user"
	"rminder/internal/pkg/logger"

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

func SetLogger(ctx *gin.Context, log *logger.Logger) {
	ctx.Set("logger", log)
}

func GetLogger(ctx *gin.Context) *logger.Logger {
	return ctx.MustGet("logger").(*logger.Logger)
}

func UserMiddleware(s *App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		user_id := user.GetUserId(session)
		if user_id == "" {
			s.logger.Warn("no user_id in session, redirecting to login", "path", ctx.Request.URL.Path)
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}

		log := s.logger.WithUserID(user_id)

		user, err := s.GetUser(user_id)
		if err != nil {
			log.Error("failed to get user, redirecting to logout", "error", err)
			ctx.Redirect(http.StatusSeeOther, "/logout")
			return
		}
		SetUser(ctx, user)

		db, err := s.GetDatabaseForUser(user_id)
		if err != nil {
			log.Error("failed to get database for user, redirecting to logout", "error", err)
			ctx.Redirect(http.StatusSeeOther, "/logout")
			return
		}
		SetUserDatabase(ctx, db)
		SetLogger(ctx, log)

		ctx.Next()
	}
}

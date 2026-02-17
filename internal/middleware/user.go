package middleware

import (
	"net/http"

	"rminder/internal/app"
	"rminder/internal/app/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func UserMiddleware(s *app.App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		user_id := user.GetUserId(session)
		if user_id == "" {
			s.Logger().Warn("no user_id in session, redirecting to login", "path", ctx.Request.URL.Path)
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}

		log := s.Logger().WithUserID(user_id).WithRequestID(app.GetRequestID(ctx))

		user, err := s.GetUser(user_id)
		if err != nil {
			log.Error("failed to get user, redirecting to logout", "error", err)
			ctx.Redirect(http.StatusSeeOther, "/logout")
			return
		}
		app.SetUser(ctx, user)

		db, err := s.GetDatabaseForUser(user_id)
		if err != nil {
			log.Error("failed to get database for user, redirecting to logout", "error", err)
			ctx.Redirect(http.StatusSeeOther, "/logout")
			return
		}
		app.SetUserDatabase(ctx, db)
		app.SetLogger(ctx, log)

		ctx.Next()
	}
}

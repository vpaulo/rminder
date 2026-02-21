package app

import (
	"net/http"
	"rminder/internal/app/user"
	"rminder/web"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LandingPageLoadHandler(s *App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		userExists := false
		user_id := user.GetUserId(session)
		if user_id != "" {
			userExists = true
		}

		ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

		err := web.Render(ctx, "home-page", map[string]any{
			"UserExists": userExists,
		})
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			s.logger.Error("error rendering in home page", "error", err)
		}
	}
}

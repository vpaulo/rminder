package app

import (
	"net/http"
	"rminder/internal/app/database"
	"rminder/internal/app/user"
	"rminder/web"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AppLoadHandler(ctx *gin.Context) {
	db := GetUserDatabase(ctx)
	log := GetLogger(ctx)

	var multiList []*database.List

	lists, err := db.Lists("")
	if err != nil {
		log.Error("error handling appLoadHandler", "error", err)
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling appLoadHandler Persistence", "error", err)
		return
	}

	// When loading lists will have the correct task count for the sidebar and
	// multiList is only for when persistence list is a list with filter so that the correct tasks are shown
	if persistence.ListId != 0 {
		filter := ""
		for _, l := range lists {
			if l.ID == persistence.ListId {
				filter = l.FilterBy
			}
		}

		if filter != "" {
			multiList, err = db.Lists(filter)
			if err != nil {
				log.Error("error handling appLoadHandler filtered lists", "error", err)
				return
			}
		}
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = web.Render(ctx.Writer, "tasks-page", map[string]any{
		"Lists":       lists,
		"MultiList":   multiList,
		"Persistence": persistence,
		"CSRFToken":   GetCSRFToken(ctx),
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		log.Error("error rendering in appLoadHandler", "error", err)
	}
}

func LandingPageLoadHandler(s *App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		userExists := false
		user_id := user.GetUserId(session)
		if user_id != "" {
			userExists = true
		}

		ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

		err := web.Render(ctx.Writer, "home-page", map[string]any{
			"UserExists": userExists,
		})
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			s.logger.Error("error rendering in home page", "error", err)
		}
	}
}

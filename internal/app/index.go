package app

import (
	"log"
	"net/http"
	"rminder/internal/app/database"
	"rminder/internal/app/user"
	"rminder/web"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AppLoadHandler(ctx *gin.Context) {
	db := GetUserDatabase(ctx)
	user := GetUser(ctx)

	var multiList []*database.List

	lists, err := db.Lists("")
	if err != nil {
		log.Fatalf("error handling appLoadHandler. Err: %v", err)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling appLoadHandler Persistence. Err: %v", err)
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
				log.Fatalf("error handling appLoadHandler filtered lists. Err: %v", err)
			}
		}
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = web.Tasks(lists, multiList, persistence, user).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("Error rendering in appLoadHandler: %e :: %v", err, e)
	}
}

func LandingPageLoadHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)

	userExists := false
	user_id := user.GetUserId(session)
	if user_id != "" {
		userExists = true
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := web.Home(userExists).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("Error rendering in Home page: %e :: %v", err, e)
	}
}

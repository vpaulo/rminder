package app

import (
	"log"
	"net/http"
	"rminder/internal/app/database"
	"rminder/web"

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
		ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("Error rendering in appLoadHandler: %e", err)
	}
}

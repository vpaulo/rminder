package routes

import (
	"log"
	"net/http"
	"rminder/internal/middleware"
	"rminder/web"

	"github.com/gin-gonic/gin"
)

func AppLoadHandler(ctx *gin.Context) {
	db := middleware.GetUserDatabase(ctx)

	lists, err := db.Lists()
	if err != nil {
		log.Fatalf("error handling appLoadHandler. Err: %v", err)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling appLoadHandler Persistence. Err: %v", err)
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = web.Tasks(lists, persistence, db.IsPremium()).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("Error rendering in appLoadHandler: %e", err)
	}
}

package app

import (
	"log"
	"net/http"
	"rminder/internal/database"
	"rminder/web"

	"github.com/gin-gonic/gin"
)

type App struct {
	db database.Service
}

func New() *App {
	return &App{
		db: database.New(),
	}
}

func (s *App) AppLoadHandler(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer

	lists, err := s.db.Lists()
	if err != nil {
		log.Fatalf("error handling appLoadHandler. Err: %v", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = web.Tasks(lists).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in appLoadHandler: %e", err)
	}
}

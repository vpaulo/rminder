package router

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	"rminder/internal/app"
	"rminder/internal/middleware"
	"rminder/internal/pkg/config"
	"rminder/internal/pkg/logger"
	"rminder/web"
)

// NewLocal registers routes without authentication, for local development use.
func NewLocal(log *logger.Logger, cfg *config.Config, dbPath string) *gin.Engine {
	application := app.New(log, cfg)

	router := gin.Default()

	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())

	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/tasks")
	})

	TasksRoutesLocal(router, application, dbPath)

	staticFiles, err := fs.Sub(web.Files, "assets")
	if err != nil {
		panic(err)
	}
	router.StaticFS("/assets", http.FS(staticFiles))

	return router
}

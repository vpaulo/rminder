package router

import (
	"encoding/gob"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"rminder/internal/app"
	"rminder/internal/login"
	"rminder/internal/login/authenticator"
	"rminder/internal/middleware"
	"rminder/internal/pkg/config"
	"rminder/internal/pkg/logger"
	"rminder/web"
)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator, log *logger.Logger, cfg *config.Config) *gin.Engine {
	application := app.New(log, cfg)

	router := gin.Default()

	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte(cfg.Auth.Cookie))
	router.Use(sessions.Sessions("auth-session", store))

	// Navigation routes
	router.GET("/login", login.LoginHandler(auth, log))
	router.GET("/callback", login.CallbackHandler(application, auth))
	router.GET("/logout", login.LogoutHandler(cfg.Auth, log))
	router.GET("/", app.LandingPageLoadHandler(application))

	// APP routes
	TasksRoutes(router, application)

	// Static files
	staticFiles, err := fs.Sub(web.Files, "assets")
	if err != nil {
		panic(err)
	}
	router.StaticFS("/assets", http.FS(staticFiles))

	return router
}

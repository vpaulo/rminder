package router

import (
	"encoding/gob"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"rminder/internal/app"
	taskhandlers "rminder/internal/handlers/tasks"
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

	router.GET("/login", login.LoginHandler(auth, log))
	router.GET("/callback", login.CallbackHandler(application, auth))
	router.GET("/logout", login.LogoutHandler(cfg.Auth, log))

	router.GET("/", app.LandingPageLoadHandler(application))
	router.GET("/tasks", middleware.UserMiddleware(application), middleware.CSRFMiddleware(), app.AppLoadHandler)

	// v0 api returns html chunks, and also to allow creation of more routes at home level
	v0 := router.Group("/v0", middleware.UserMiddleware(application), middleware.CSRFMiddleware())

	tasks := v0.Group("/tasks")
	tasks.GET("/all", taskhandlers.GetTasks)
	tasks.GET("/my-day", taskhandlers.GetTasks)
	tasks.GET("/important", taskhandlers.GetTasks)
	tasks.GET("/completed", taskhandlers.GetTasks)
	tasks.POST("/create", taskhandlers.CreateTask)
	tasks.GET("/:taskID", taskhandlers.GetTask)
	tasks.DELETE("/:taskID", taskhandlers.DeleteTask)
	tasks.PUT("/:taskID/:slug", taskhandlers.UpdateTask)

	lists := v0.Group("/lists")
	lists.GET("/all", taskhandlers.GetLists)
	lists.POST("/create", taskhandlers.CreateList)
	lists.POST("/search", taskhandlers.SearchLists)
	lists.GET("/:listID", taskhandlers.GetList)
	lists.DELETE("/:listID", taskhandlers.DeleteList)
	lists.PUT("/:listID", taskhandlers.UpdateList)

	// v1 api returns JSON
	SetV1Routes(router, application)

	// Static files
	staticFiles, err := fs.Sub(web.Files, "assets")
	if err != nil {
		panic(err)
	}
	router.StaticFS("/assets", http.FS(staticFiles))

	return router
}

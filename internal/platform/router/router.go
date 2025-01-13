package router

import (
	"encoding/gob"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"rminder/internal/app"
	"rminder/internal/app/callback"
	"rminder/internal/app/login"
	"rminder/internal/app/logout"
	"rminder/internal/platform/authenticator"
	"rminder/internal/platform/middleware"
	"rminder/web"
)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator) *gin.Engine {
	app := app.New()

	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	router.GET("/login", login.Handler(auth))
	router.GET("/callback", callback.Handler(auth))
	router.GET("/logout", logout.Handler)

	router.GET("/", middleware.IsAuthenticated, app.AppLoadHandler)

	tasks := router.Group("/tasks", middleware.IsAuthenticated)
	tasks.GET("/all", app.GetTasks)
	tasks.GET("/my-day", app.GetTasks)
	tasks.GET("/important", app.GetTasks)
	tasks.GET("/completed", app.GetTasks)
	tasks.POST("/create", app.CreateTask)
	tasks.GET("/:taskID", app.GetTask)
	tasks.DELETE("/:taskID", app.DeleteTask)
	tasks.PUT("/:taskID/:slug", app.UpdateTask)

	lists := router.Group("/lists", middleware.IsAuthenticated)
	lists.GET("/all", app.GetTasks)
	lists.POST("/create", app.CreateList)
	lists.GET("/:listID", app.GetList)
	lists.DELETE("/:listID", app.DeleteList)
	lists.PUT("/:listID/:slug", app.UpdateList)

	// Static files
	staticFiles, err := fs.Sub(web.Files, "assets")
	if err != nil {
		panic(err)
	}
	router.StaticFS("/assets", http.FS(staticFiles))

	return router
}

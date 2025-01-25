package router

import (
	"encoding/gob"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"rminder/internal/app"
	"rminder/internal/authenticator"
	"rminder/internal/middleware"
	"rminder/internal/routes"
	"rminder/web"
)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator) *gin.Engine {
	application := app.New()

	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	router.GET("/login", routes.LoginHandler(auth))
	router.GET("/callback", routes.CallbackHandler(application, auth))
	router.GET("/logout", routes.LogoutHandler)

	router.GET("/", middleware.Authentication(application), routes.AppLoadHandler)

	checkout := router.Group("/checkout", middleware.Authentication(application))
	checkout.POST("/create-checkout-session", routes.CreatePremiumCheckoutSession)
	checkout.GET("/success", routes.PremiumCheckoutSuccessHandler)

	router.POST("/post-checkout/webhook", routes.CheckoutWebhookHandler(application))

	tasks := router.Group("/tasks", middleware.Authentication(application))
	tasks.GET("/all", routes.GetTasks)
	tasks.GET("/my-day", routes.GetTasks)
	tasks.GET("/important", routes.GetTasks)
	tasks.GET("/completed", routes.GetTasks)
	tasks.POST("/create", routes.CreateTask)
	tasks.GET("/:taskID", routes.GetTask)
	tasks.DELETE("/:taskID", routes.DeleteTask)
	tasks.PUT("/:taskID/:slug", routes.UpdateTask)

	lists := router.Group("/lists", middleware.Authentication(application))
	lists.GET("/all", routes.GetLists)
	lists.POST("/create", routes.CreateList)
	lists.POST("/search", routes.SearchLists)
	lists.GET("/:listID", routes.GetList)
	lists.DELETE("/:listID", routes.DeleteList)
	lists.PUT("/:listID", routes.UpdateList)

	// Static files
	staticFiles, err := fs.Sub(web.Files, "assets")
	if err != nil {
		panic(err)
	}
	router.StaticFS("/assets", http.FS(staticFiles))

	return router
}

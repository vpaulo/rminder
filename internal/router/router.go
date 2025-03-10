package router

import (
	"encoding/gob"
	"io/fs"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"rminder/internal/app"
	"rminder/internal/checkout"
	"rminder/internal/login"
	"rminder/internal/login/authenticator"
	"rminder/web"
)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator) *gin.Engine {
	application := app.New()

	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte(os.Getenv("COOKIE_AUTHENTICATION_KEY")))
	router.Use(sessions.Sessions("auth-session", store))

	router.GET("/login", login.LoginHandler(auth))
	router.GET("/callback", login.CallbackHandler(application, auth))
	router.GET("/logout", login.LogoutHandler)

	router.GET("/", app.LandingPageLoadHandler)
	router.GET("/tasks", app.UserMiddleware(application), app.AppLoadHandler)

	checkout_group := router.Group("/checkout", app.UserMiddleware(application))
	checkout_group.POST("/create-checkout-session", checkout.CreatePremiumCheckoutSession)
	checkout_group.GET("/success", checkout.PremiumCheckoutSuccessHandler)

	router.POST("/post-checkout/webhook", checkout.CheckoutWebhookHandler(application))

	// v0 api returns html chunks, and also to allow creation of more routes at home level
	v0 := router.Group("/v0", app.UserMiddleware(application))

	tasks := v0.Group("/tasks")
	tasks.GET("/all", app.GetTasks)
	tasks.GET("/my-day", app.GetTasks)
	tasks.GET("/important", app.GetTasks)
	tasks.GET("/completed", app.GetTasks)
	tasks.POST("/create", app.CreateTask)
	tasks.GET("/:taskID", app.GetTask)
	tasks.DELETE("/:taskID", app.DeleteTask)
	tasks.PUT("/:taskID/:slug", app.UpdateTask)

	lists := v0.Group("/lists")
	lists.GET("/all", app.GetLists)
	lists.POST("/create", app.CreateList)
	lists.POST("/search", app.SearchLists)
	lists.GET("/:listID", app.GetList)
	lists.DELETE("/:listID", app.DeleteList)
	lists.PUT("/:listID", app.UpdateList)

	// Static files
	staticFiles, err := fs.Sub(web.Files, "assets")
	if err != nil {
		panic(err)
	}
	router.StaticFS("/assets", http.FS(staticFiles))

	return router
}

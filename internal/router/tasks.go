package router

import (
	"rminder/internal/app"
	taskhandlers "rminder/internal/handlers/tasks"
	"rminder/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TasksRoutes(router *gin.Engine, application *app.App) {
	// Navigation routes
	router.GET("/tasks", middleware.UserMiddleware(application), middleware.CSRFMiddleware(), taskhandlers.Load)

	// Partials: returns HTML chunks
	partials := router.Group("/partials", middleware.UserMiddleware(application), middleware.CSRFMiddleware())

	tasks := partials.Group("/tasks")
	tasks.GET("/all", taskhandlers.GetTasks)
	tasks.GET("/my-day", taskhandlers.GetTasks)
	tasks.GET("/important", taskhandlers.GetTasks)
	tasks.GET("/completed", taskhandlers.GetTasks)
	tasks.POST("/create", taskhandlers.CreateTask)
	tasks.GET("/:taskID", taskhandlers.GetTask)
	tasks.DELETE("/:taskID", taskhandlers.DeleteTask)
	tasks.PUT("/:taskID/:slug", taskhandlers.UpdateTask)

	lists := partials.Group("/lists")
	lists.GET("/all", taskhandlers.GetLists)
	lists.POST("/create", taskhandlers.CreateList)
	lists.POST("/search", taskhandlers.SearchLists)
	lists.GET("/:listID", taskhandlers.GetList)
	lists.DELETE("/:listID", taskhandlers.DeleteList)
	lists.PUT("/:listID", taskhandlers.UpdateList)

	// API: returns JSON
	api := router.Group("/api", middleware.UserMiddleware(application), middleware.CSRFMiddleware())

	apiTasks := api.Group("/tasks")
	apiTasks.GET("/export", taskhandlers.ExportLists)
	apiTasks.POST("/import", taskhandlers.ImportLists)
	apiTasks.POST("/reorder", taskhandlers.ReorderTasks)

	apiLists := api.Group("/lists")
	apiLists.POST("/reorder", taskhandlers.ReorderLists)
}

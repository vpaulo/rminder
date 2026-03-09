package router

import (
	"rminder/internal/app"
	taskhandlers "rminder/internal/handlers/tasks"
	"rminder/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TasksRoutesLocal(router *gin.Engine, application *app.App, dbPath string) {
	local := middleware.LocalMiddleware(dbPath, application)

	router.GET("/tasks", local, taskhandlers.Load)

	partials := router.Group("/partials", local)

	tasks := partials.Group("/tasks")
	tasks.GET("/all", taskhandlers.GetTasks)
	tasks.GET("/my-day", taskhandlers.GetTasks)
	tasks.GET("/important", taskhandlers.GetTasks)
	tasks.GET("/completed", taskhandlers.GetTasks)
	tasks.POST("/create", taskhandlers.CreateTask)
	tasks.GET("/:taskID", taskhandlers.GetTask)
	tasks.DELETE("/:taskID", taskhandlers.DeleteTask)
	tasks.PUT("/:taskID/:slug", taskhandlers.UpdateTask)
	tasks.POST("/:taskID/subtask", taskhandlers.CreateSubtask)

	lists := partials.Group("/lists")
	lists.GET("/all", taskhandlers.GetLists)
	lists.POST("/create", taskhandlers.CreateList)
	lists.POST("/search", taskhandlers.SearchLists)
	lists.GET("/:listID", taskhandlers.GetList)
	lists.DELETE("/:listID", taskhandlers.DeleteList)
	lists.PUT("/:listID", taskhandlers.UpdateList)

	api := router.Group("/api", local)

	apiTasks := api.Group("/tasks")
	apiTasks.GET("/export", taskhandlers.ExportLists)
	apiTasks.POST("/import", taskhandlers.ImportLists)
	apiTasks.POST("/reorder", taskhandlers.ReorderTasks)

	apiLists := api.Group("/lists")
	apiLists.POST("/reorder", taskhandlers.ReorderLists)
}

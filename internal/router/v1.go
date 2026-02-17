package router

import (
	"rminder/internal/app"
	taskhandlers "rminder/internal/handlers/tasks"
	"rminder/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetV1Routes(router *gin.Engine, application *app.App) {
	// v1 api returns JSON
	v1 := router.Group("/v1", middleware.UserMiddleware(application), middleware.CSRFMiddleware())

	v1.GET("/export", taskhandlers.ExportLists)
	v1.POST("/import", taskhandlers.ImportLists)

	tasks := v1.Group("/tasks")
	tasks.POST("/reorder", taskhandlers.ReorderTasks)

	lists := v1.Group("/lists")
	lists.POST("/reorder", taskhandlers.ReorderLists)
}

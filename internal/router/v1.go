package router

import (
	"rminder/internal/app"

	"github.com/gin-gonic/gin"
)

func SetV1Routes(router *gin.Engine, application *app.App) {
	// v1 api returns JSON
	v1 := router.Group("/v1", app.UserMiddleware(application))

	tasks := v1.Group("/tasks")
	tasks.POST("/reorder", app.ReorderTasks)

	lists := v1.Group("/lists")
	lists.POST("/reorder", app.ReorderLists)
}

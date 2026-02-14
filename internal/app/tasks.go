package app

import (
	"errors"
	"net/http"
	"rminder/internal/app/database"
	"rminder/web"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetTasks(ctx *gin.Context) {
	db := GetUserDatabase(ctx)
	log := GetLogger(ctx)

	slug := strings.TrimPrefix(ctx.Request.URL.Path, "/")

	var (
		tasks []*database.Task
		lists []*database.List
		err   error
	)

	switch slug {
	case "important":
		tasks, err = db.ImportantTasks()
	case "completed":
		tasks, err = db.CompletedTasks()
	default:
		lists, err = db.Lists("")
	}

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Error("error handling tasks", "error", err)
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling GetTasks Persistence", "error", err)
		return
	}

	if slug == "" {
		ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = web.Render(ctx.Writer, "tasks-page", map[string]any{
			"Lists":       lists,
			"MultiList":   ([]*database.List)(nil),
			"Persistence": persistence,
		})
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Error("error rendering in tasksHandler", "error", err)
		}
	} else {
		ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = web.Render(ctx.Writer, "task-list", map[string]any{
			"Tasks":        tasks,
			"SelectedTask": persistence.TaskId,
		})
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Error("error rendering in TaskList", "error", err)
		}
	}
}

func GetTask(ctx *gin.Context) {
	db := GetUserDatabase(ctx)
	log := GetLogger(ctx)

	taskID := ctx.Param("taskID")

	var (
		task *database.Task
		err  error
	)

	if taskID != "" {
		task, err = db.Task(taskID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Error("error handling task", "taskID", taskID, "error", err)
			return
		}

		id, _ := strconv.Atoi(taskID)
		err = db.UpdatePersistenceTask(id)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Error("error updating persistence task", "taskID", taskID, "error", err)
			return
		}
	} else {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("task ID must not be empty"))
		log.Error("error task id is empty")
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.Render(ctx.Writer, "details", task)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Error("error rendering in TaskList", "error", err)
	}
}

func CreateTask(ctx *gin.Context) {
	db := GetUserDatabase(ctx)
	log := GetLogger(ctx)

	err := ctx.Request.ParseForm()
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		log.Error("error parsing form", "error", err)
		return
	}

	title := ctx.Request.FormValue("task")
	list, e := strconv.Atoi(ctx.Request.FormValue("list"))

	if e != nil {
		list = 0
	}

	if title != "" && len(title) >= 3 && len(title) <= 255 && list != 0 {
		err := db.CreateTask(title, list)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Error("error creating task", "error", err)
			return
		}
		log.Info("task created", "title", title, "listID", list)
	} else {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("title validation failed"))
		log.Error("error title validation failed", "title", title, "listID", list)
		return
	}

	tasks, err := db.ListTasks(list, "")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Error("error handling tasks", "error", err)
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling CreateTask Persistence", "error", err)
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.Render(ctx.Writer, "task-list", map[string]any{
		"Tasks":        tasks,
		"SelectedTask": persistence.TaskId,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Error("error rendering in TaskList", "error", err)
	}
}

func DeleteTask(ctx *gin.Context) {
	db := GetUserDatabase(ctx)
	log := GetLogger(ctx)

	taskID := ctx.Param("taskID")

	var err error

	if taskID != "" {
		err = db.DeleteTask(taskID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Error("error deleting task", "taskID", taskID, "error", err)
			return
		}
		log.Info("task deleted", "taskID", taskID)
	} else {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("task ID must not be empty"))
		log.Error("error task id is empty")
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = ctx.Writer.Write([]byte(""))
}

func UpdateTask(ctx *gin.Context) {
	db := GetUserDatabase(ctx)
	log := GetLogger(ctx)

	taskID := ctx.Param("taskID")
	slug := ctx.Param("slug")

	err := ctx.Request.ParseForm()
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		log.Error("error parsing form", "error", err)
		return
	}

	if taskID != "" {
		var err error

		switch slug {
		case "title":
			title := ctx.Request.FormValue("title")
			if title != "" && len(title) >= 3 && len(title) <= 255 {
				err = db.UpdateTask(taskID, title)
			} else {
				ctx.AbortWithError(http.StatusBadRequest, errors.New("title validation failed"))
				log.Error("error title validation failed", "taskID", taskID)
				return
			}
		case "description":
			err = db.UpdateTaskDescription(taskID, ctx.Request.FormValue("description"))
		case "important":
			err = db.ToggleImportant(taskID)
		case "completed":
			err = db.ToggleComplete(taskID)
		case "priority":
			err = db.UpdateTaskPriority(taskID, ctx.Request.FormValue("priority"))
		case "date-start":
			err = db.UpdateTaskStartDate(taskID, ctx.Request.FormValue("from"))
		case "date-end":
			err = db.UpdateTaskEndDate(taskID, ctx.Request.FormValue("to"))
		case "remove-persistence":
			err = db.UpdatePersistenceTask(0)
		}

		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Error("error updating task", "taskID", taskID, "slug", slug, "error", err)
			return
		}
		log.Info("task updated", "taskID", taskID, "field", slug)
	} else {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("task ID must not be empty"))
		log.Error("error task id is empty")
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling UpdateTask Persistence", "error", err)
		return
	}

	if slug == "description" {
		// TODO return proper message
		_, _ = ctx.Writer.Write([]byte("Updated description"))
	} else {
		// get task
		task, err := db.Task(taskID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Error("error handling task", "taskID", taskID, "error", err)
			return
		}

		ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

		if slug == "important" {
			err = web.Render(ctx.Writer, "task-important-elem", task)
		} else if slug == "completed" {
			err = web.Render(ctx.Writer, "task-completed-elem", task)
		} else {
			err = web.Render(ctx.Writer, "task", map[string]any{
				"Task":         task,
				"SelectedTask": persistence.TaskId,
			})
		}

		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Error("error rendering in Task", "error", err)
		}
	}
}

func ReorderTasks(ctx *gin.Context) {
	db := GetUserDatabase(ctx)
	log := GetLogger(ctx)
	var reorder []database.Reorder
	var err error

	if err = ctx.ShouldBindJSON(&reorder); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		log.Error("error binding reorder JSON", "error", err)
		return
	}

	err = db.ReorderTasks(reorder)

	if err != nil {
		log.Error("error reordering tasks", "error", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Tasks order update unsuccessful.",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	log.Info("tasks reordered")
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message": "Tasks order update successful.",
		"status":  http.StatusOK,
	})
}

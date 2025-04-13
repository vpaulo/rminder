package app

import (
	"errors"
	"log"
	"net/http"
	"rminder/internal/app/database"
	"rminder/web"
	"rminder/web/components"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetTasks(ctx *gin.Context) {
	db := GetUserDatabase(ctx)
	user := GetUser(ctx)

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
		e := ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("error handling tasks. Err: %e :: %v", err, e)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling GetTasks Persistence. Err: %v", err)
	}

	if slug == "" {
		ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = web.Tasks(lists, nil, persistence, user).Render(ctx.Request.Context(), ctx.Writer)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("Error rendering in tasksHandler: %e :: %v", err, e)
		}
	} else {
		ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = web.TaskList(tasks, persistence.TaskId).Render(ctx.Request.Context(), ctx.Writer)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("Error rendering in TaskList: %e :: %v", err, e)
		}
	}
}

func GetTask(ctx *gin.Context) {
	db := GetUserDatabase(ctx)

	taskID := ctx.Param("taskID")

	var (
		task *database.Task
		err  error
	)

	if taskID != "" {
		task, err = db.Task(taskID)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error handling task. Err: %e :: %v", err, e)
		}

		id, _ := strconv.Atoi(taskID)
		err = db.UpdatePersistenceTask(id)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error updating persistence task. Err: %e :: %v", err, e)
		}
	} else {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error task id. Err: %e :: %v", err, e)
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = components.Details(task).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		e := ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("Error rendering in TaskList: %e :: %v", err, e)
	}
}

func CreateTask(ctx *gin.Context) {
	db := GetUserDatabase(ctx)

	err := ctx.Request.ParseForm()
	if err != nil {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error parsing form. Err: %e :: %v", err, e)
	}

	// create new task
	title := ctx.Request.FormValue("task")
	list, e := strconv.Atoi(ctx.Request.FormValue("list"))

	if e != nil {
		list = 0
	}

	if title != "" && len(title) >= 3 && len(title) <= 255 && list != 0 {
		err := db.CreateTask(title, list)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error creating task. Err: %e :: %v", err, e)
		}
	} else {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error title validation failed. Err: %e :: %v", err, e)
	}

	tasks, err := db.ListTasks(list, "")
	if err != nil {
		e := ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("error handling tasks. Err: %e :: %v", err, e)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling CreateTask Persistence. Err: %v", err)
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.TaskList(tasks, persistence.TaskId).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		e := ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("Error rendering in TaskList: %e :: %v", err, e)
	}
}

func DeleteTask(ctx *gin.Context) {
	db := GetUserDatabase(ctx)

	taskID := ctx.Param("taskID")

	var err error

	if taskID != "" {
		err = db.DeleteTask(taskID)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error deleting task. Err: %e :: %v", err, e)
		}
	} else {
		e := ctx.AbortWithError(http.StatusBadRequest, errors.New("task ID must not be empty"))
		log.Fatalf("error task id. Err: %e :: %v", err, e)
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = ctx.Writer.Write([]byte(""))
}

func UpdateTask(ctx *gin.Context) {
	db := GetUserDatabase(ctx)

	taskID := ctx.Param("taskID")
	slug := ctx.Param("slug")

	err := ctx.Request.ParseForm()
	if err != nil {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error parsing form. Err: %e :: %v", err, e)
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
				e := ctx.AbortWithError(http.StatusInternalServerError, errors.New("title validation failed"))
				log.Fatalf("error title validation failed. Err: %e :: %v", err, e)
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
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error updating task. Err: %e :: %v", err, e)
		}
	} else {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error task id. Err: %e :: %v", err, e)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling UpdateTask Persistence. Err: %e", err)
	}

	if slug == "description" {
		// TODO return proper message
		_, _ = ctx.Writer.Write([]byte("Updated description"))
	} else {
		// get task
		task, err := db.Task(taskID)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error handling task. Err: %e :: %v", err, e)
		}

		ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

		if slug == "important" {
			err = web.TaskImportantElem(task).Render(ctx.Request.Context(), ctx.Writer)
		} else if slug == "completed" {
			err = web.TaskCompletedElem(task).Render(ctx.Request.Context(), ctx.Writer)
		} else {
			err = web.Task(task, persistence.TaskId).Render(ctx.Request.Context(), ctx.Writer)
		}

		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("Error rendering in Task: %e :: %v", err, e)
		}
	}
}

func ReorderTasks(ctx *gin.Context) {
	db := GetUserDatabase(ctx)
	var reorder []database.Reorder
	var err error

	if err = ctx.ShouldBindJSON(&reorder); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// reorder tasks.
	err = db.ReorderTasks(reorder)

	if err != nil {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"message": "Tasks order update unsuccessful.",
			"status":  http.StatusInternalServerError,
		})
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message": "Tasks order update successful.",
		"status":  http.StatusOK,
	})
}

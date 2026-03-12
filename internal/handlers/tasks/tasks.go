package tasks

import (
	"net/http"
	"rminder/internal/app"
	"rminder/internal/app/database"
	"rminder/internal/app/sse"
	"rminder/web"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func publishEvent(ctx *gin.Context, eventType string) {
	broker := app.GetSSEBroker(ctx)
	user := app.GetUser(ctx)
	broker.Publish(user.Id, sse.Event{Type: eventType, Time: time.Now().UTC().Format(time.RFC3339)})
}

func Load(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	var multiList []*database.List

	lists, err := db.Lists("")
	if err != nil {
		log.Error("error handling Tasks app Load", "error", err)
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling Tasks app Load Persistence", "error", err)
		return
	}

	// When loading lists will have the correct task count for the sidebar and
	// multiList is only for when persistence list is a list with filter so that the correct tasks are shown
	if persistence.ListId != 0 {
		filter := ""
		for _, l := range lists {
			if l.ID == persistence.ListId {
				filter = l.FilterBy
			}
		}

		if filter != "" {
			multiList, err = db.Lists(filter)
			if err != nil {
				log.Error("error handling Tasks app Load filtered lists", "error", err)
				return
			}
		}
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = web.Render(ctx.Writer, "tasks-page", map[string]any{
		"Lists":       lists,
		"MultiList":   multiList,
		"Persistence": persistence,
		"CSRFToken":   app.GetCSRFToken(ctx),
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		log.Error("error rendering in Tasks app Load", "error", err)
	}
}

func GetTasks(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

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
		log.Error("error handling tasks", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load tasks.")
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling GetTasks Persistence", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load tasks.")
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
			log.Error("error rendering in tasksHandler", "error", err)
			app.ErrorInternalHTML(ctx, "Failed to render tasks.")
		}
	} else {
		ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = web.Render(ctx.Writer, "task-list", map[string]any{
			"Tasks":        tasks,
			"SelectedTask": persistence.TaskId,
		})
		if err != nil {
			log.Error("error rendering in TaskList", "error", err)
			app.ErrorInternalHTML(ctx, "Failed to render task list.")
		}
	}
}

func GetTask(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	taskID := ctx.Param("taskID")

	var (
		task *database.Task
		err  error
	)

	if taskID != "" {
		task, err = db.Task(taskID)
		if err != nil {
			log.Error("error handling task", "taskID", taskID, "error", err)
			app.ErrorInternalHTML(ctx, "Failed to load task.")
			return
		}

		id, _ := strconv.Atoi(taskID)
		err = db.UpdatePersistenceTask(id)
		if err != nil {
			log.Error("error updating persistence task", "taskID", taskID, "error", err)
			app.ErrorInternalHTML(ctx, "Failed to update task state.")
			return
		}
	} else {
		log.Error("error task id is empty")
		app.ErrorBadRequestHTML(ctx, "Task ID must not be empty.")
		return
	}

	subtasks, err := db.Subtasks(task.ID)
	if err != nil {
		log.Error("error fetching subtasks", "taskID", taskID, "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load subtasks.")
		return
	}

	allLists, err := db.Lists("")
	if err != nil {
		log.Error("error fetching lists", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load lists.")
		return
	}
	var realLists []*database.List
	for _, l := range allLists {
		if l.FilterBy == "" {
			realLists = append(realLists, l)
		}
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.Render(ctx.Writer, "details", map[string]any{
		"Task":     task,
		"Subtasks": subtasks,
		"Lists":    realLists,
	})
	if err != nil {
		log.Error("error rendering in TaskList", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to render task details.")
	}
}

func CreateTask(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	err := ctx.Request.ParseForm()
	if err != nil {
		log.Error("error parsing form", "error", err)
		app.ErrorBadRequestHTML(ctx, "Invalid form data.")
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
			log.Error("error creating task", "error", err)
			app.ErrorInternalHTML(ctx, "Failed to create task.")
			return
		}
		publishEvent(ctx, "task_created")
		log.Info("task created", "title", title, "listID", list)
	} else {
		log.Error("error title validation failed", "title", title, "listID", list)
		app.ErrorBadRequestHTML(ctx, "Title must be between 3 and 255 characters and a list must be selected.")
		return
	}

	tasks, err := db.ListTasks(list, "")
	if err != nil {
		log.Error("error handling tasks", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load tasks.")
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling CreateTask Persistence", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load task state.")
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.Render(ctx.Writer, "task-list", map[string]any{
		"Tasks":        tasks,
		"SelectedTask": persistence.TaskId,
	})
	if err != nil {
		log.Error("error rendering in TaskList", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to render task list.")
	}
}

func DeleteTask(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	taskID := ctx.Param("taskID")

	var err error

	if taskID != "" {
		// Delete subtasks first, then the parent
		id, _ := strconv.Atoi(taskID)
		subtasks, _ := db.Subtasks(id)
		for _, sub := range subtasks {
			_ = db.DeleteTask(strconv.Itoa(sub.ID))
		}

		err = db.DeleteTask(taskID)
		if err != nil {
			log.Error("error deleting task", "taskID", taskID, "error", err)
			app.ErrorInternalHTML(ctx, "Failed to delete task.")
			return
		}
		publishEvent(ctx, "task_deleted")
		log.Info("task deleted", "taskID", taskID)
	} else {
		log.Error("error task id is empty")
		app.ErrorBadRequestHTML(ctx, "Task ID must not be empty.")
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = ctx.Writer.Write([]byte(""))
}

func UpdateTask(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	taskID := ctx.Param("taskID")
	slug := ctx.Param("slug")

	err := ctx.Request.ParseForm()
	if err != nil {
		log.Error("error parsing form", "error", err)
		app.ErrorBadRequestHTML(ctx, "Invalid form data.")
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
				log.Error("error title validation failed", "taskID", taskID)
				app.ErrorBadRequestHTML(ctx, "Title must be between 3 and 255 characters.")
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
		case "list":
			err = db.UpdateTaskList(taskID, ctx.Request.FormValue("list"))
		case "date-start":
			err = db.UpdateTaskStartDate(taskID, ctx.Request.FormValue("from"))
		case "date-end":
			err = db.UpdateTaskEndDate(taskID, ctx.Request.FormValue("to"))
		case "remove-persistence":
			err = db.UpdatePersistenceTask(0)
		}

		if err != nil {
			log.Error("error updating task", "taskID", taskID, "slug", slug, "error", err)
			app.ErrorInternalHTML(ctx, "Failed to update task.")
			return
		}
		publishEvent(ctx, "task_updated")
		log.Info("task updated", "taskID", taskID, "field", slug)
	} else {
		log.Error("error task id is empty")
		app.ErrorBadRequestHTML(ctx, "Task ID must not be empty.")
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling UpdateTask Persistence", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load task state.")
		return
	}

	if slug == "description" {
		// TODO return proper message
		_, _ = ctx.Writer.Write([]byte("Updated description"))
	} else if slug == "list" {
		_, _ = ctx.Writer.Write([]byte(""))
	} else {
		// get task
		task, err := db.Task(taskID)
		if err != nil {
			log.Error("error handling task", "taskID", taskID, "error", err)
			app.ErrorInternalHTML(ctx, "Failed to load task.")
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
			log.Error("error rendering in Task", "error", err)
			app.ErrorInternalHTML(ctx, "Failed to render task.")
		}
	}
}

func CreateSubtask(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	parentID := ctx.Param("taskID")

	err := ctx.Request.ParseForm()
	if err != nil {
		log.Error("error parsing form", "error", err)
		app.ErrorBadRequestHTML(ctx, "Invalid form data.")
		return
	}

	title := ctx.Request.FormValue("task")
	pid, e := strconv.Atoi(parentID)
	if e != nil || pid == 0 {
		log.Error("invalid parent task id", "parentID", parentID)
		app.ErrorBadRequestHTML(ctx, "Invalid parent task ID.")
		return
	}

	parent, err := db.Task(parentID)
	if err != nil {
		log.Error("error fetching parent task", "parentID", parentID, "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load parent task.")
		return
	}

	if title != "" && len(title) >= 3 && len(title) <= 255 {
		err = db.CreateSubtask(title, parent.ListId, pid)
		if err != nil {
			log.Error("error creating subtask", "error", err)
			app.ErrorInternalHTML(ctx, "Failed to create subtask.")
			return
		}
		publishEvent(ctx, "subtask_created")
		log.Info("subtask created", "title", title, "parentID", pid)
	} else {
		log.Error("error title validation failed", "title", title)
		app.ErrorBadRequestHTML(ctx, "Title must be between 3 and 255 characters.")
		return
	}

	subtasks, err := db.Subtasks(pid)
	if err != nil {
		log.Error("error fetching subtasks", "parentID", pid, "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load subtasks.")
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.Render(ctx.Writer, "subtask-list", map[string]any{
		"Subtasks": subtasks,
		"ParentID": pid,
	})
	if err != nil {
		log.Error("error rendering subtask-list", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to render subtasks.")
	}
}

func ReorderTasks(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)
	var reorder []database.Reorder
	var err error

	if err = ctx.ShouldBindJSON(&reorder); err != nil {
		log.Error("error binding reorder JSON", "error", err)
		app.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request data.")
		return
	}

	err = db.ReorderTasks(reorder)

	if err != nil {
		log.Error("error reordering tasks", "error", err)
		app.ErrorJSON(ctx, http.StatusInternalServerError, "Tasks order update unsuccessful.")
		return
	}

	log.Info("tasks reordered")
	publishEvent(ctx, "task_reordered")
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message": "Tasks order update successful.",
		"status":  http.StatusOK,
	})
}

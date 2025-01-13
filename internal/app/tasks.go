package app

import (
	"log"
	"net/http"
	"rminder/internal/database"
	"rminder/web"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *App) GetTasks(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer

	slug := strings.TrimPrefix(r.URL.Path, "/")

	var (
		tasks []*database.Task
		lists []*database.List
		err   error
	)

	switch slug {
	case "my-day":
		// tasks, err = s.db.MyDayTasks()
	case "important":
		tasks, err = s.db.ImportantTasks()
	case "completed":
		tasks, err = s.db.CompletedTasks()
	default:
		lists, err = s.db.Lists()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("error handling tasks. Err: %v", err)
	}

	if slug == "" {
		// TODO: find better way to update totals of tasks lists
		// totals, err := s.db.Totals()
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	log.Fatalf("error handling totals. Err: %v", err)
		// }

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = web.Tasks(lists).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("Error rendering in tasksHandler: %e", err)
		}
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = web.TaskList(tasks).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("Error rendering in TaskList: %e", err)
		}
	}
}

func (s *App) GetTask(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer

	taskID := ctx.Param("taskID")

	var (
		task *database.Task
		err  error
	)

	if taskID != "" {
		task, err = s.db.Task(taskID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error handling task. Err: %v", err)
		}
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.TaskDetails(task).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Error rendering in TaskList: %e", err)
	}
}

func (s *App) CreateTask(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Fatalf("error parsing form. Err: %v", err)
	}

	// create new task
	title := r.FormValue("task")
	list, e := strconv.Atoi(r.FormValue("list"))

	if e != nil {
		list = 0
	}

	if title != "" && len(title) >= 3 && len(title) <= 255 && list != 0 {
		err := s.db.CreateTask(title, list)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error creating task. Err: %v", err)
		}
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("error title validation failed. Err: %v", err)
	}

	tasks, err := s.db.ListTasks(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("error handling tasks. Err: %v", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.TaskList(tasks).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Error rendering in TaskList: %e", err)
	}
}

func (s *App) DeleteTask(ctx *gin.Context) {
	w := ctx.Writer

	taskID := ctx.Param("taskID")

	var err error

	if taskID != "" {
		err = s.db.DeleteTask(taskID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error deleting task. Err: %v", err)
		}
	} else {
		http.Error(w, "Task ID must not be empty", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(""))
}

func (s *App) UpdateTask(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer

	taskID := ctx.Param("taskID")
	slug := ctx.Param("slug")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	if taskID != "" {
		var err error

		switch slug {
		case "title":
			title := r.FormValue("title")
			if title != "" && len(title) >= 3 && len(title) <= 255 {
				err = s.db.UpdateTask(taskID, title)
			} else {
				// TODO use ApiError here
				http.Error(w, "Title validation failed", http.StatusInternalServerError)
				log.Fatalf("error title validation failed. Err: %v", err)
			}
		case "description":
			err = s.db.UpdateTaskDescription(taskID, r.FormValue("description"))
		case "important":
			err = s.db.ToggleImportant(taskID)
		case "completed":
			err = s.db.ToggleComplete(taskID)
		case "my-day":
			// err = s.db.ToggleMyDay(taskID)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error updating task. Err: %v", err)
		}
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if slug == "description" {
		// TODO return proper message
		_, _ = w.Write([]byte("Updated description"))
	} else {
		// get task
		task, err := s.db.Task(taskID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error handling task. Err: %v", err)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = web.Task(task).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("Error rendering in Task: %e", err)
		}
	}
}

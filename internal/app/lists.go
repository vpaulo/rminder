package app

import (
	"encoding/json"
	"log"
	"net/http"
	"rminder/internal/database"
	"rminder/web"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *App) GetLists(ctx *gin.Context) {
	w := ctx.Writer

	lists, err := s.db.Lists()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("error handling lists. Err: %v", err)
	}

	jsonResp, err := json.Marshal(lists)

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *App) GetList(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer

	listID := ctx.Param("listID")

	var (
		list *database.List
		err  error
	)

	if listID != "" {
		list, err = s.db.List(listID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error handling task. Err: %v", err)
		}

		id, _ := strconv.Atoi(listID)
		err = s.db.UpdatePersistence(0, id, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error updating persistence list. Err: %v", err)
		}
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	persistence, err := s.db.Persistence()
	if err != nil {
		log.Fatalf("error handling GetList Persistence. Err: %v", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.ListsContent(list, persistence).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Error rendering in ListsContent: %e", err)
	}
}

func (s *App) CreateList(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Fatalf("error parsing form. Err: %v", err)
	}

	// create new task
	list := r.FormValue("new-list")

	pos, e := strconv.Atoi(r.FormValue("position"))

	if e != nil {
		pos = 0
	}

	swatch := r.FormValue("swatch")
	icon := r.FormValue("icon")
	if list != "" && len(list) >= 3 && len(list) <= 255 && pos != 0 && swatch != "" && icon != "" {
		err := s.db.CreateList(list, swatch, icon, pos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error creating list. Err: %v", err)
		}
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("error form field validation failed. Err: %v", err)
	}

	lists, err := s.db.Lists()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("error handling lists. Err: %v", err)
	}

	persistence, err := s.db.Persistence()
	if err != nil {
		log.Fatalf("error handling CreateList Persistence. Err: %v", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.Lists(lists, false, persistence.ListId).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Error rendering in Lists: %e", err)
	}
}

func (s *App) DeleteList(ctx *gin.Context) {
	w := ctx.Writer

	listID := ctx.Param("listID")

	var err error

	if listID != "" {
		err = s.db.DeleteTask(listID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error deleting task. Err: %v", err)
		}
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(""))
}

// TODO: update this to work with lists
func (s *App) UpdateList(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer

	listID := ctx.Param("listID")
	slug := r.PathValue("slug")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	if listID != "" {
		var err error

		switch slug {
		case "title":
			title := r.FormValue("title")
			if title != "" && len(title) >= 3 && len(title) <= 255 {
				err = s.db.UpdateTask(listID, title)
			} else {
				// TODO use ApiError here
				http.Error(w, "Title validation failed", http.StatusInternalServerError)
				log.Fatalf("error title validation failed. Err: %v", err)
			}
		case "description":
			err = s.db.UpdateTaskDescription(listID, r.FormValue("description"))
		case "important":
			err = s.db.ToggleImportant(listID)
		case "completed":
			err = s.db.ToggleComplete(listID)
		case "my-day":
			// err = s.db.ToggleMyDay(listID)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error updating task. Err: %v", err)
		}
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	persistence, err := s.db.Persistence()
	if err != nil {
		log.Fatalf("error handling UpdateList Persistence. Err: %v", err)
	}

	if slug == "description" {
		// TODO return proper message
		_, _ = w.Write([]byte("Updated description"))
	} else {
		// get task
		task, err := s.db.Task(listID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error handling task. Err: %v", err)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = web.Task(task, persistence.TaskId).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("Error rendering in Task: %e", err)
		}
	}
}

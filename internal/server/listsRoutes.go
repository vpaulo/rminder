package server

import (
	"encoding/json"
	"log"
	"net/http"
	"rminder/internal/database"
	"rminder/web"
)

func (s *Server) ListsRoutes() *http.ServeMux {
	routes := http.NewServeMux()

	routes.HandleFunc("GET /all", s.getTasks)
	routes.HandleFunc("POST /create", s.createList)
	routes.HandleFunc("GET /{listID}", s.getList)
	routes.HandleFunc("DELETE /{listID}", s.deleteList)
	routes.HandleFunc("PUT /{listID}/{slug}", s.updateList)

	return routes
}

func (s *Server) getLists(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) getList(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskID")

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

func (s *Server) createList(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Fatalf("error parsing form. Err: %v", err)
	}

	// create new task
	title := r.FormValue("task")
	if title != "" && len(title) >= 3 && len(title) <= 255 {
		err := s.db.CreateTask(title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error creating task. Err: %v", err)
		}
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("error title validation failed. Err: %v", err)
	}

	tasks, err := s.db.Tasks()
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

func (s *Server) deleteList(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskID")

	var err error

	if taskID != "" {
		err = s.db.DeleteTask(taskID)
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

func (s *Server) updateList(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskID")
	slug := r.PathValue("slug")

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

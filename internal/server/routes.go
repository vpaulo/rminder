package server

import (
	"log"
	"net/http"

	mw "rminder/internal/middleware"
	"rminder/web"
)

func (s *Server) RegisterRoutes() http.Handler {
	stack := mw.MiddlewareStack(
		mw.Logger,
	)

	router := http.NewServeMux()
	router.Handle("/tasks/", http.StripPrefix("/tasks", s.TasksRoutes()))

	// Static files
	fileServer := http.FileServer(http.FS(web.Files))
	router.Handle("GET /assets/", fileServer)

	// Views/pages
	router.HandleFunc("/{$}", s.tasksHandler)

	return stack(router)
}

func (s *Server) tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.db.Tasks()
	if err != nil {
		log.Fatalf("error handling tasks. Err: %v", err)
	}

	totals, err := s.db.Totals()
	if err != nil {
		log.Fatalf("error handling totals. Err: %v", err)
	}

	err = web.Tasks(tasks, totals).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in tasksHandler: %e", err)
	}
}

package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"rminder/cmd/web"
)

type Middleware func(http.Handler) http.Handler

func (s *Server) RegisterRoutes() http.Handler {
	stack := MiddlewareStack(
		Logger,
		// IsLoggedIn,
	)

	router := http.NewServeMux()

	apiRouter := http.NewServeMux()
	apiRouter.HandleFunc("GET /task", s.getAllTasks)
	apiRouter.HandleFunc("POST /task", s.createTask)
	apiRouter.HandleFunc("GET /task/{taskID}", s.getTask)
	apiRouter.HandleFunc("PUT /task/{taskID}", s.updateTask)
	apiRouter.HandleFunc("DELETE /task/{taskID}", s.deleteTask)
	router.Handle("/api/v0/", http.StripPrefix("/api/v0", apiRouter))

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

func (s *Server) getAllTasks(w http.ResponseWriter, r *http.Request) {
	// TODO: instead of log.Fatalf maybe send message to FE that something went wrong or just log and return zero tasks
	tasks, err := s.db.Tasks()
	if err != nil {
		log.Fatalf("error handling tasks. Err: %v", err)
	}

	jsonResp, err := json.Marshal(tasks)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	// TODO: create tasks from db
	fmt.Fprintln(w, "Create task")
}

func (s *Server) getTask(w http.ResponseWriter, r *http.Request) {
	// TODO: get task from db
	fmt.Fprintln(w, "Get task by ID")
}

func (s *Server) updateTask(w http.ResponseWriter, r *http.Request) {
	// TODO: update task from db
	fmt.Fprintln(w, "Update task")
}

func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	// TODO: delete task from db
	fmt.Fprintln(w, "Delete task")
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

// MiddlewareStack chains multiple middlewares
func MiddlewareStack(ms ...Middleware) Middleware {
	return Middleware(func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			m := ms[i]
			next = m(next)
		}
		return next
	})
}

func IsCustomer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("checking if is api...")
		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(time.Since(start), r.Method, r.URL.Path)
	})
}

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

type ApiError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("Api error: %d", e.StatusCode)
}

func NewApiError(statusCode int, error error) ApiError {
	return ApiError{
		StatusCode: statusCode,
		Msg:        error.Error(),
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	stack := MiddlewareStack(
		Logger,
		// IsLoggedIn,
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

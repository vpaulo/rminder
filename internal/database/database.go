package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	Tasks() ([]*Task, error)
	MyDayTasks() ([]*Task, error)
	ImportantTasks() ([]*Task, error)
	CompletedTasks() ([]*Task, error)
	Totals() (*Total, error)
	CreateTask(title string) error
	ToggleComplete(ID string) error
	ToggleImportant(ID string) error
	ToggleMyDay(ID string) error
}

type service struct {
	db *sql.DB
}

var (
	dburl      = os.Getenv("DB_URL")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("sqlite", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	dbInstance = &service{
		db: db,
	}

	err = dbInstance.loadSqlFile()
	if err != nil {
		// Failed to create tables.
		log.Fatal(err)
	}

	return dbInstance
}

func (s *service) loadSqlFile() error {
	// Read file
	file, err := Schema.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	// Execute all
	_, err = s.db.Exec(string(file))
	if err != nil {
		return err
	}

	return nil
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	return s.db.Close()
}

func (s *service) Tasks() ([]*Task, error) {
	query, err := s.db.Prepare("SELECT * FROM task ORDER BY updated_at DESC")
	if err != nil {
		return nil, fmt.Errorf("DB.Tasks - prepare query failed: %v", err)
	}
	defer query.Close()

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("DB.Tasks - query result failed: %v", err)
	}

	tasks := make([]*Task, 0)
	for result.Next() {
		data := new(Task)
		err := result.Scan(
			&data.ID,
			&data.Title,
			&data.Description,
			&data.Completed,
			&data.Important,
			&data.MyDay,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.Tasks - result scan failed: %v", err)
		}
		tasks = append(tasks, data)
	}

	return tasks, nil
}

func (s *service) MyDayTasks() ([]*Task, error) {
	query, err := s.db.Prepare("SELECT * FROM task WHERE my_day = true ORDER BY updated_at DESC")
	if err != nil {
		return nil, fmt.Errorf("DB.MyDayTasks - prepare query failed: %v", err)
	}
	defer query.Close()

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("DB.MyDayTasks - query result failed: %v", err)
	}

	tasks := make([]*Task, 0)
	for result.Next() {
		data := new(Task)
		err := result.Scan(
			&data.ID,
			&data.Title,
			&data.Description,
			&data.Completed,
			&data.Important,
			&data.MyDay,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.MyDayTasks - result scan failed: %v", err)
		}
		tasks = append(tasks, data)
	}

	return tasks, nil
}

func (s *service) ImportantTasks() ([]*Task, error) {
	query, err := s.db.Prepare("SELECT * FROM task WHERE important = true ORDER BY updated_at DESC")
	if err != nil {
		return nil, fmt.Errorf("DB.ImportantTasks - prepare query failed: %v", err)
	}
	defer query.Close()

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("DB.ImportantTasks - query result failed: %v", err)
	}

	tasks := make([]*Task, 0)
	for result.Next() {
		data := new(Task)
		err := result.Scan(
			&data.ID,
			&data.Title,
			&data.Description,
			&data.Completed,
			&data.Important,
			&data.MyDay,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.ImportantTasks - result scan failed: %v", err)
		}
		tasks = append(tasks, data)
	}

	return tasks, nil
}

func (s *service) CompletedTasks() ([]*Task, error) {
	query, err := s.db.Prepare("SELECT * FROM task WHERE completed = true ORDER BY updated_at DESC")
	if err != nil {
		return nil, fmt.Errorf("DB.CompletedTasks - prepare query failed: %v", err)
	}
	defer query.Close()

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("DB.CompletedTasks - query result failed: %v", err)
	}

	tasks := make([]*Task, 0)
	for result.Next() {
		data := new(Task)
		err := result.Scan(
			&data.ID,
			&data.Title,
			&data.Description,
			&data.Completed,
			&data.Important,
			&data.MyDay,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.CompletedTasks - result scan failed: %v", err)
		}
		tasks = append(tasks, data)
	}

	return tasks, nil
}

func (s *service) Totals() (*Total, error) {
	query, err := s.db.Prepare(`SELECT
		    COALESCE(SUM(CASE WHEN completed = true THEN 1 ELSE 0 END), 0) AS total_completed,
			COALESCE(SUM(CASE WHEN important = true THEN 1 ELSE 0 END), 0) AS total_important,
			COALESCE(SUM(CASE WHEN my_day = true THEN 1 ELSE 0 END), 0) AS total_my_day,
			COUNT(*) as total_tasks	FROM task`)
	if err != nil {
		return nil, fmt.Errorf("DB.Totals - prepare query failed: %v", err)
	}
	defer query.Close()

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("DB.Totals - query result failed: %v", err)
	}

	total := make([]*Total, 0)
	for result.Next() {
		data := new(Total)
		err := result.Scan(
			&data.Completed,
			&data.Important,
			&data.MyDay,
			&data.Tasks,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.Totals - result scan failed: %v", err)
		}
		total = append(total, data)
	}

	// Return the only row
	return total[0], nil
}

func (s *service) CreateTask(title string) error {
	query, err := s.db.Prepare("INSERT INTO task (title) Values (?)")
	if err != nil {
		return fmt.Errorf("DB.CreateTask - prepare create query failed: %v", err)
	}
	defer query.Close()

	task := &Task{
		Title: title,
	}

	_, err = query.Exec(task.Title)
	if err != nil {
		return fmt.Errorf("DB.CreateTask - create query result failed: %v", err)
	}

	return nil
}

func (s *service) ToggleComplete(ID string) error {
	query, err := s.db.Prepare("UPDATE task SET completed = NOT completed, updated_at = CURRENT_TIMESTAMP WHERE task_id=?")
	if err != nil {
		return fmt.Errorf("DB.ToggleComplete - prepare update query failed: %v", err)
	}
	defer query.Close()

	_, err = query.Exec(ID)
	if err != nil {
		return fmt.Errorf("DB.ToggleComplete - update query result failed: %v", err)
	}

	return nil
}

func (s *service) ToggleImportant(ID string) error {
	query, err := s.db.Prepare("UPDATE task SET important = NOT important, updated_at = CURRENT_TIMESTAMP WHERE task_id=?")
	if err != nil {
		return fmt.Errorf("DB.ToggleImportant - prepare update query failed: %v", err)
	}
	defer query.Close()

	_, err = query.Exec(ID)
	if err != nil {
		return fmt.Errorf("DB.ToggleImportant - update query result failed: %v", err)
	}

	return nil
}

func (s *service) ToggleMyDay(ID string) error {
	query, err := s.db.Prepare("UPDATE task SET my_day = NOT my_day, updated_at = CURRENT_TIMESTAMP WHERE task_id=?")
	if err != nil {
		return fmt.Errorf("DB.ToggleMyDay - prepare update query failed: %v", err)
	}
	defer query.Close()

	_, err = query.Exec(ID)
	if err != nil {
		return fmt.Errorf("DB.ToggleMyDay - update query result failed: %v", err)
	}

	return nil
}

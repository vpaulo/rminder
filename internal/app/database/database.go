package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	ImportantTasks() ([]*Task, error)
	CompletedTasks() ([]*Task, error)
	CreateTask(title string, list int) error
	ToggleComplete(ID string) error
	ToggleImportant(ID string) error
	Task(ID string) (*Task, error)
	UpdateTask(ID string, title string) error
	UpdateTaskDescription(ID string, description string) error
	DeleteTask(ID string) error
	UpdateTaskPriority(ID string, priority string) error
	UpdateTaskStartDate(id string, date string) error
	UpdateTaskEndDate(id string, date string) error

	Lists() ([]*List, error)
	List(ID string) (*List, error)
	ListTasks(id int) ([]*Task, error)
	CreateList(name string, swatch string, icon string, position int, pinned bool, filter string) error
	UpdateList(id int, name string, colour string, icon string, pinned bool, filter string) error
	DeleteList(id int) error

	Groups() ([]*GroupList, error)
	Group(ID int) (*GroupList, error)
	GroupLists(id int) ([]*List, error)

	Persistence() (*Persistence, error)
	UpdatePersistence(task int, list int, group int) error
	UpdatePersistenceTask(task int) error
	UpdatePersistenceList(list int) error
	UpdatePersistenceGroup(group int) error

	SearchLists(searchQuery string) ([]*List, error)
	ListTasksSearch(id int, searchQuery string) ([]*Task, error)
}

type service struct {
	database_path string
	db            *sql.DB
}

func New(database_path string) Service {
	db, err := sql.Open("sqlite", database_path)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatalf("%s: %v", database_path, err)
	}

	dbInstance := &service{
		database_path: database_path,
		db:            db,
	}

	err = dbInstance.loadSqlFile()
	if err != nil {
		// Failed to create tables.
		log.Fatalf("%s: %v", database_path, err)
	}

	return dbInstance
}

func (s *service) loadSqlFile() error {
	// Check db already has been initialised
	lists, err := s.Lists()
	if len(lists) > 0 || err == nil {
		return nil
	}

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
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
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
	log.Printf("Disconnected from database: %s", s.database_path)
	return s.db.Close()
}

func (s *service) Tasks() ([]*Task, error) {
	// TODO: maybe change the default for order by created_at
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
			&data.Priority,
			&data.Position,
			&data.StartAt,
			&data.EndAt,
			&data.ListId,
			&data.ParentId,
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
			&data.Priority,
			&data.Position,
			&data.StartAt,
			&data.EndAt,
			&data.ListId,
			&data.ParentId,
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
			&data.Priority,
			&data.Position,
			&data.StartAt,
			&data.EndAt,
			&data.ListId,
			&data.ParentId,
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

// func (s *service) Totals() (*Total, error) {
// 	query, err := s.db.Prepare(`SELECT
// 		    COALESCE(SUM(CASE WHEN completed = true THEN 1 ELSE 0 END), 0) AS total_completed,
// 			COALESCE(SUM(CASE WHEN important = true THEN 1 ELSE 0 END), 0) AS total_important,
// 			COALESCE(SUM(CASE WHEN my_day = true THEN 1 ELSE 0 END), 0) AS total_my_day,
// 			COUNT(*) as total_tasks	FROM task`)
// 	defer query.Close()
// 	if err != nil {
// 		return nil, fmt.Errorf("DB.Totals - prepare query failed: %v", err)
// 	}

// 	total := new(Total)
// 	err = query.QueryRow().Scan(
// 		&total.Completed,
// 		&total.Important,
// 		&total.MyDay,
// 		&total.Tasks)
// 	if err != nil {
// 		return nil, fmt.Errorf("DB.Totals - query result failed: %v", err)
// 	}

// 	return total, nil
// }

func (s *service) CreateTask(title string, list int) error {
	query, err := s.db.Prepare("INSERT INTO task (title, list_id) Values (?,?)")
	if err != nil {
		return fmt.Errorf("DB.CreateTask - prepare create query failed: %v", err)
	}
	defer query.Close()

	task := &Task{
		Title:  title,
		ListId: list,
	}

	_, err = query.Exec(task.Title, task.ListId)
	if err != nil {
		return fmt.Errorf("DB.CreateTask - create query result failed: %v", err)
	}

	return nil
}

func (s *service) ToggleComplete(id string) error {
	query, err := s.db.Prepare("UPDATE task SET completed = NOT completed, updated_at = CURRENT_TIMESTAMP WHERE id=?")
	if err != nil {
		return fmt.Errorf("DB.ToggleComplete - prepare update query failed: %v", err)
	}
	defer query.Close()

	_, err = query.Exec(id)
	if err != nil {
		return fmt.Errorf("DB.ToggleComplete - update query result failed: %v", err)
	}

	return nil
}

func (s *service) ToggleImportant(id string) error {
	query, err := s.db.Prepare("UPDATE task SET important = NOT important, updated_at = CURRENT_TIMESTAMP WHERE id=?")
	if err != nil {
		return fmt.Errorf("DB.ToggleImportant - prepare update query failed: %v", err)
	}
	defer query.Close()

	_, err = query.Exec(id)
	if err != nil {
		return fmt.Errorf("DB.ToggleImportant - update query result failed: %v", err)
	}

	return nil
}

func (s *service) Task(id string) (*Task, error) {
	query, err := s.db.Prepare("SELECT * FROM task WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("DB.Task - prepare query failed: %v", err)
	}
	defer query.Close()

	task := new(Task)
	err = query.QueryRow(id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.Important,
		&task.Priority,
		&task.Position,
		&task.StartAt,
		&task.EndAt,
		&task.ListId,
		&task.ParentId,
		&task.CreatedAt,
		&task.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("DB.Task - query result failed: %v", err)
	}

	return task, nil
}

func (s *service) UpdateTask(id string, title string) error {
	query, err := s.db.Prepare("UPDATE task SET title = ?, updated_at = CURRENT_TIMESTAMP WHERE id=?")
	if err != nil {
		return fmt.Errorf("DB.UpdateTask - prepare update query failed: %v", err)
	}
	defer query.Close()

	_, err = query.Exec(title, id)
	if err != nil {
		return fmt.Errorf("DB.UpdateTask - update query result failed: %v", err)
	}

	return nil
}

func (s *service) UpdateTaskDescription(id string, description string) error {
	query, err := s.db.Prepare("UPDATE task SET description = ?, updated_at = CURRENT_TIMESTAMP WHERE id=?")
	if err != nil {
		return fmt.Errorf("DB.UpdateTaskDescription - prepare update query failed: %v", err)
	}
	defer query.Close()

	_, err = query.Exec(description, id)
	if err != nil {
		return fmt.Errorf("DB.UpdateTaskDescription - update query result failed: %v", err)
	}

	return nil
}

func (s *service) UpdateTaskPriority(id string, priority string) error {
	query, err := s.db.Prepare("UPDATE task SET priority = ?, updated_at = CURRENT_TIMESTAMP WHERE id=?")
	if err != nil {
		return fmt.Errorf("DB.UpdateTaskPriority - prepare update query failed: %v", err)
	}
	defer query.Close()

	_, err = query.Exec(priority, id)
	if err != nil {
		return fmt.Errorf("DB.UpdateTaskPriority - update query result failed: %v", err)
	}

	return nil
}

func (s *service) UpdateTaskStartDate(id string, date string) error {
	query, err := s.db.Prepare("UPDATE task SET start_at = ?, updated_at = CURRENT_TIMESTAMP WHERE id=?")
	if err != nil {
		return fmt.Errorf("DB.UpdateTaskStartDate - prepare update query failed: %v", err)
	}
	defer query.Close()

	// TODO: create date formating helpers to be used across application
	tm, err := time.Parse("2006-01-02", date)

	if err != nil {
		return fmt.Errorf("DB.UpdateTaskStartDate - format date failed: %v", err)
	}

	_, err = query.Exec(tm.Format(time.DateTime), id)
	if err != nil {
		return fmt.Errorf("DB.UpdateTaskStartDate - update query result failed: %v", err)
	}

	return nil
}

func (s *service) UpdateTaskEndDate(id string, date string) error {
	query, err := s.db.Prepare("UPDATE task SET end_at = ?, updated_at = CURRENT_TIMESTAMP WHERE id=?")
	if err != nil {
		return fmt.Errorf("DB.UpdateTaskEndDate - prepare update query failed: %v", err)
	}
	defer query.Close()

	// TODO: create date formating helpers to be used across application
	tm, err := time.Parse("2006-01-02", date)

	if err != nil {
		return fmt.Errorf("DB.UpdateTaskEndDate - format date failed: %v", err)
	}

	_, err = query.Exec(tm.Format(time.DateTime), id)
	if err != nil {
		return fmt.Errorf("DB.UpdateTaskEndDate - update query result failed: %v", err)
	}

	return nil
}

func (s *service) DeleteTask(id string) error {
	query, err := s.db.Prepare("DELETE FROM task WHERE id=?")
	if err != nil {
		return fmt.Errorf("DB.DeleteTask - prepare update query failed: %v", err)
	}
	defer query.Close()

	_, err = query.Exec(id)
	if err != nil {
		return fmt.Errorf("DB.DeleteTask - update query result failed: %v", err)
	}

	return nil
}

func (s *service) Lists() ([]*List, error) {
	query, err := s.db.Prepare("SELECT * FROM list ORDER BY created_at ASC")
	if err != nil {
		return nil, fmt.Errorf("DB.Lists - prepare query failed: %v", err)
	}
	defer query.Close()

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("DB.Lists - query result failed: %v", err)
	}

	var tasks []*Task

	lists := make([]*List, 0)
	for result.Next() {
		data := new(List)
		err := result.Scan(
			&data.ID,
			&data.Name,
			&data.Colour,
			&data.Icon,
			&data.FilterBy,
			&data.GroupId,
			&data.Pinned,
			&data.Base,
			&data.Position,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.Lists - result scan failed: %v", err)
		}
		tasks, err = s.ListTasks(data.ID)
		if err != nil {
			return nil, fmt.Errorf("DB.Lists - get list tasks failed: %v", err)
		}
		data.Tasks = tasks
		lists = append(lists, data)
	}

	return lists, nil
}

func (s *service) ListTasks(id int) ([]*Task, error) {
	query, err := s.db.Prepare("SELECT * FROM task WHERE list_id=?")
	if err != nil {
		return nil, fmt.Errorf("DB.ListTasks - prepare query failed: %v", err)
	}
	defer query.Close()

	result, err := query.Query(id)
	if err != nil {
		return nil, fmt.Errorf("DB.ListTasks - query result failed: %v", err)
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
			&data.Priority,
			&data.Position,
			&data.StartAt,
			&data.EndAt,
			&data.ListId,
			&data.ParentId,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.ListTasks - result scan failed: %v", err)
		}
		tasks = append(tasks, data)
	}

	return tasks, nil
}

func (s *service) List(id string) (*List, error) {
	query, err := s.db.Prepare("SELECT * FROM list WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("DB.List - prepare query failed: %v", err)
	}
	defer query.Close()

	var tasks []*Task

	list := new(List)
	err = query.QueryRow(id).Scan(
		&list.ID,
		&list.Name,
		&list.Colour,
		&list.Icon,
		&list.FilterBy,
		&list.GroupId,
		&list.Pinned,
		&list.Base,
		&list.Position,
		&list.CreatedAt,
		&list.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("DB.List - query result failed: %v", err)
	}

	tasks, err = s.ListTasks(list.ID)
	if err != nil {
		return nil, fmt.Errorf("DB.List - get list tasks failed: %v", err)
	}
	list.Tasks = tasks

	return list, nil
}

func (s *service) Groups() ([]*GroupList, error) {
	query, err := s.db.Prepare("SELECT * FROM group_list ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("DB.Groups - prepare query failed: %v", err)
	}
	defer query.Close()

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("DB.Groups - query result failed: %v", err)
	}

	var lists []*List

	groups := make([]*GroupList, 0)
	for result.Next() {
		data := new(GroupList)
		err := result.Scan(
			&data.ID,
			&data.Name,
			&data.Position,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.Groups - result scan failed: %v", err)
		}
		lists, err = s.GroupLists(data.ID)
		if err != nil {
			return nil, fmt.Errorf("DB.Groups - get list tasks failed: %v", err)
		}
		data.Lists = lists
		groups = append(groups, data)
	}

	return groups, nil
}

func (s *service) Group(id int) (*GroupList, error) {
	query, err := s.db.Prepare("SELECT * FROM group_list WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("DB.Group - prepare query failed: %v", err)
	}
	defer query.Close()

	var lists []*List

	group := new(GroupList)
	err = query.QueryRow(id).Scan(
		&group.ID,
		&group.Name,
		&group.Position,
		&group.CreatedAt,
		&group.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("DB.Group - query result failed: %v", err)
	}

	lists, err = s.GroupLists(group.ID)
	if err != nil {
		return nil, fmt.Errorf("DB.Group - get list tasks failed: %v", err)
	}
	group.Lists = lists

	return group, nil
}

func (s *service) GroupLists(id int) ([]*List, error) {
	query, err := s.db.Prepare("SELECT * FROM list WHERE group_id=?")
	if err != nil {
		return nil, fmt.Errorf("DB.GroupLists - prepare query failed: %v", err)
	}
	defer query.Close()

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("DB.GroupLists - query result failed: %v", err)
	}

	var tasks []*Task

	lists := make([]*List, 0)
	for result.Next() {
		data := new(List)
		err := result.Scan(
			&data.ID,
			&data.Name,
			&data.Colour,
			&data.Icon,
			&data.FilterBy,
			&data.GroupId,
			&data.Pinned,
			&data.Base,
			&data.Position,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.GroupLists - result scan failed: %v", err)
		}
		tasks, err = s.ListTasks(data.ID)
		if err != nil {
			return nil, fmt.Errorf("DB.GroupLists - get list tasks failed: %v", err)
		}
		data.Tasks = tasks
		lists = append(lists, data)
	}

	return lists, nil
}

func (s *service) CreateList(name string, swatch string, icon string, position int, pinned bool, filter string) error {
	query, err := s.db.Prepare("INSERT INTO list (name, colour, icon, position, pinned, filter_by) Values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("DB.CreateList - prepare create query failed: %v", err)
	}
	defer query.Close()

	list := &List{
		Name:     name,
		Colour:   swatch,
		Icon:     icon,
		Position: position,
		Pinned:   pinned,
		FilterBy: filter,
	}

	_, err = query.Exec(list.Name, list.Colour, list.Icon, list.Position, list.Pinned, list.FilterBy)
	if err != nil {
		return fmt.Errorf("DB.CreateList - create query result failed: %v", err)
	}

	return nil
}

func (s *service) UpdateList(id int, name string, colour string, icon string, pinned bool, filter string) error {
	query, err := s.db.Prepare("UPDATE list SET name = ?, colour = ?, icon = ?, pinned = ?, filter_by = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?")
	if err != nil {
		return fmt.Errorf("DB.UpdateList - prepare update query failed: %v", err)
	}
	defer query.Close()

	_, err = query.Exec(name, colour, icon, pinned, filter, id)
	if err != nil {
		return fmt.Errorf("DB.UpdateList - update query result failed: %v", err)
	}

	return nil
}

func (s *service) DeleteList(id int) error {
	query, err := s.db.Prepare("DELETE FROM list WHERE id=?")
	if err != nil {
		return fmt.Errorf("DB.DeleteList - prepare update query failed: %v", err)
	}
	defer query.Close()

	_, err = query.Exec(id)
	if err != nil {
		return fmt.Errorf("DB.DeleteList - update query result failed: %v", err)
	}

	return nil
}

func (s *service) Persistence() (*Persistence, error) {
	query, err := s.db.Prepare("SELECT * FROM persistence WHERE id=1")
	defer query.Close()
	if err != nil {
		return nil, fmt.Errorf("DB.Persistence - prepare query failed: %v", err)
	}

	data := new(Persistence)
	err = query.QueryRow().Scan(
		&data.ID,
		&data.TaskId,
		&data.ListId,
		&data.GroupId)

	if err != nil {
		return nil, fmt.Errorf("DB.Persistence - query result failed: %v", err)
	}

	return data, nil
}

func (s *service) UpdatePersistence(task int, list int, group int) error {
	query, err := s.db.Prepare("UPDATE persistence SET task_id = ?, list_id = ?, group_id = ? WHERE id=1")
	defer query.Close()
	if err != nil {
		return fmt.Errorf("DB.UpdatePersistence - prepare update query failed: %v", err)
	}

	_, err = query.Exec(task, list, group)
	if err != nil {
		return fmt.Errorf("DB.UpdatePersistence - update query result failed: %v", err)
	}

	return nil
}

func (s *service) UpdatePersistenceTask(task int) error {
	query, err := s.db.Prepare("UPDATE persistence SET task_id = ? WHERE id=1")
	defer query.Close()
	if err != nil {
		return fmt.Errorf("DB.UpdatePersistenceTask - prepare update query failed: %v", err)
	}

	_, err = query.Exec(task)
	if err != nil {
		return fmt.Errorf("DB.UpdatePersistenceTask - update query result failed: %v", err)
	}

	return nil
}

func (s *service) UpdatePersistenceList(list int) error {
	query, err := s.db.Prepare("UPDATE persistence SET list_id = ? WHERE id=1")
	defer query.Close()
	if err != nil {
		return fmt.Errorf("DB.UpdatePersistenceList - prepare update query failed: %v", err)
	}

	_, err = query.Exec(list)
	if err != nil {
		return fmt.Errorf("DB.UpdatePersistenceList - update query result failed: %v", err)
	}

	return nil
}

func (s *service) UpdatePersistenceGroup(group int) error {
	query, err := s.db.Prepare("UPDATE persistence SET group_id = ? WHERE id=1")
	defer query.Close()
	if err != nil {
		return fmt.Errorf("DB.UpdatePersistenceGroup - prepare update query failed: %v", err)
	}

	_, err = query.Exec(group)
	if err != nil {
		return fmt.Errorf("DB.UpdatePersistenceGroup - update query result failed: %v", err)
	}

	return nil
}

func (s *service) SearchLists(searchQuery string) ([]*List, error) {
	query, err := s.db.Prepare("SELECT * FROM list ORDER BY created_at ASC")
	if err != nil {
		return nil, fmt.Errorf("DB.SearchLists - prepare query failed: %v", err)
	}
	defer query.Close()

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("DB.SearchLists - query result failed: %v", err)
	}

	var tasks []*Task

	lists := make([]*List, 0)
	for result.Next() {
		data := new(List)
		err := result.Scan(
			&data.ID,
			&data.Name,
			&data.Colour,
			&data.Icon,
			&data.FilterBy,
			&data.GroupId,
			&data.Pinned,
			&data.Base,
			&data.Position,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.SearchLists - result scan failed: %v", err)
		}
		tasks, err = s.ListTasksSearch(data.ID, searchQuery)
		if err != nil {
			return nil, fmt.Errorf("DB.SearchLists - get list tasks failed: %v", err)
		}
		data.Tasks = tasks

		if len(tasks) > 0 {
			lists = append(lists, data)
		}
	}

	return lists, nil
}

func (s *service) ListTasksSearch(id int, searchQuery string) ([]*Task, error) {
	query, err := s.db.Prepare("SELECT * FROM task WHERE list_id=? AND (title || description) LIKE ?")
	if err != nil {
		return nil, fmt.Errorf("DB.ListTasksSearch - prepare query failed: %v", err)
	}
	defer query.Close()

	result, err := query.Query(id, searchQuery+"%")
	if err != nil {
		return nil, fmt.Errorf("DB.ListTasksSearch - query result failed: %v", err)
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
			&data.Priority,
			&data.Position,
			&data.StartAt,
			&data.EndAt,
			&data.ListId,
			&data.ParentId,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.ListTasksSearch - result scan failed: %v", err)
		}
		tasks = append(tasks, data)
	}

	return tasks, nil
}

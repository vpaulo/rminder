package database

import "embed"

//go:embed "migrations"
var Migrations embed.FS

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Important   bool   `json:"important"`
	Priority    int    `json:"priority"`
	Position    int    `json:"position"`
	StartAt     string `json:"start_at"`
	EndAt       string `json:"end_at"`
	ListId      int    `json:"list_id"`
	ParentId    int    `json:"parent_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type List struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Colour    string `json:"colour"`
	Icon      string `json:"icon"`
	FilterBy  string `json:"filter_by"`
	Pinned    bool   `json:"pinned"`
	Base      bool   `json:"base"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Tasks     []*Task
}

type Total struct {
	Tasks     int `json:"total_tasks"`
	Completed int `json:"total_completed"`
	Important int `json:"total_important"`
}

type Persistence struct {
	ID     int `json:"id"`
	TaskId int `json:"task_id"`
	ListId int `json:"list_id"`
}

type Reorder struct {
	ID       int `json:"id"`
	Position int `json:"position"`
}

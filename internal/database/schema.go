package database

import "embed"

//go:embed "schema.sql"
var Schema embed.FS

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
	GroupId   int    `json:"group_id"`
	Pinned    bool   `json:"pinned"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Tasks     []*Task
}

type GroupList struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Lists     []*List
}

type Total struct {
	Tasks     int `json:"total_tasks"`
	Completed int `json:"total_completed"`
	Important int `json:"total_important"`
}

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

type Total struct {
	Tasks     int `json:"total_tasks"`
	Completed int `json:"total_completed"`
	Important int `json:"total_important"`
}

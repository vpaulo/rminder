package database

import "embed"

//go:embed "schema.sql"
var Schema embed.FS

type Task struct {
	ID          int    `json:"task_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Important   bool   `json:"important"`
	MyDay       bool   `json:"my_day"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Total struct {
	Tasks     int `json:"total_tasks"`
	Completed int `json:"total_completed"`
	Important int `json:"total_important"`
	MyDay     int `json:"total_my_day"`
}

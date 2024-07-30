package database

import "embed"

//go:embed "schema.sql"
var Schema embed.FS

type Task struct {
	ID          int    `json:"task_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

package app

import (
	"rminder/internal/database"
)

type App struct {
	db database.Service
}

func New() *App {
	return &App{
		db: database.New(),
	}
}

func (s *App) GetDatabaseForUser(user_id string) (database.Service, error) {
	// TODO: implement user database cache
	return s.db, nil
}

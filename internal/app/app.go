package app

import (
	"os"
	"path"
	"rminder/internal/database"
)

type App struct {
	user_databases map[string]database.Service
}

func New() *App {
	return &App{
		user_databases: make(map[string]database.Service),
	}
}

func ensureDirectoryExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	} else if err != nil {
		return err
	}
	return nil
}

func (s *App) GetDatabaseForUser(user_id string) (database.Service, error) {
	if db, ok := s.user_databases[user_id]; !ok {
		user_database_root_directory := os.Getenv("USER_DATABASE_ROOT_DIRECTORY")
		user_database_directory := path.Join(user_database_root_directory, user_id)
		ensureDirectoryExists(user_database_directory)

		user_database_path := path.Join(user_database_directory, "db.sqlite")
		db := database.New(user_database_path)
		s.user_databases[user_id] = db
		return db, nil
	} else {
		return db, nil
	}
}

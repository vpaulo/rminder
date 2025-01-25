package app

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"rminder/internal/app/database"
	"rminder/internal/app/user"
)

type App struct {
	user_databases map[string]database.Service
	users          map[string]*user.User
}

func New() *App {
	return &App{
		user_databases: make(map[string]database.Service),
		users:          make(map[string]*user.User),
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
		err := ensureDirectoryExists(user_database_directory)
		if err != nil {
			log.Fatalf("Failed to create directory: %s error: %v", user_database_directory, err)
		}

		user_database_path := path.Join(user_database_directory, "db.sqlite")
		db := database.New(user_database_path)
		s.user_databases[user_id] = db
		return db, nil
	} else {
		return db, nil
	}
}

func (s *App) loadUserFromFile(user_id string) (*user.User, error) {
	user_database_root_directory := os.Getenv("USER_DATABASE_ROOT_DIRECTORY")
	user_database_directory := path.Join(user_database_root_directory, user_id)
	user_file_path := path.Join(user_database_directory, "user.json")

	file, err := os.Open(user_file_path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var user user.User
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *App) GetUser(user_id string) (*user.User, error) {
	if user_obj, ok := s.users[user_id]; !ok {
		user_obj, err := s.loadUserFromFile(user_id)
		if err != nil {
			return nil, err
		}
		s.users[user_id] = user_obj
		return user_obj, nil
	} else {
		return user_obj, nil
	}
}

func (s *App) SaveUser(user *user.User) error {
	user_database_root_directory := os.Getenv("USER_DATABASE_ROOT_DIRECTORY")
	user_database_directory := path.Join(user_database_root_directory, user.Id)
	user_file_path := path.Join(user_database_directory, "user.json")

	err := ensureDirectoryExists(user_database_directory)
	if err != nil {
		return err
	}

	file, err := os.Create(user_file_path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(user)
	if err != nil {
		return err
	}

	return nil
}

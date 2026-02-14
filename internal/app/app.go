package app

import (
	"encoding/json"
	"os"
	"path"

	"rminder/internal/app/database"
	"rminder/internal/app/user"
	"rminder/internal/pkg/config"
	"rminder/internal/pkg/logger"
)

type App struct {
	user_databases map[string]database.Service
	users          map[string]*user.User
	logger         *logger.Logger
	config         *config.Config
}

func New(log *logger.Logger, cfg *config.Config) *App {
	return &App{
		user_databases: make(map[string]database.Service),
		users:          make(map[string]*user.User),
		logger:         log,
		config:         cfg,
	}
}

func (s *App) Logger() *logger.Logger {
	return s.logger
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
		user_directory := path.Join(s.config.Database.UsersDir, user_id)
		err := ensureDirectoryExists(user_directory)
		if err != nil {
			s.logger.Error("Failed to create directory: %s error: %v", user_directory, err)
			return nil, err
		}

		user_database_path := path.Join(user_directory, "db.sqlite")
		db := database.New(user_database_path)
		s.user_databases[user_id] = db
		return db, nil
	} else {
		return db, nil
	}
}

func (s *App) loadUserFromFile(user_id string) (*user.User, error) {
	user_directory := path.Join(s.config.Database.UsersDir, user_id)
	user_file_path := path.Join(user_directory, "user.json")

	file, err := os.Open(user_file_path)
	if err != nil {
		s.logger.Error("Failed to load user from file: %s error: %v", user_file_path, err)
		return nil, err
	}
	defer file.Close()

	var user user.User
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&user)
	if err != nil {
		s.logger.Error("Failed to decode error: %v", "error", err)
		return nil, err
	}

	return &user, nil
}

func (s *App) GetUser(user_id string) (*user.User, error) {
	if user_obj, ok := s.users[user_id]; !ok {
		user_obj, err := s.loadUserFromFile(user_id)
		if err != nil {
			s.logger.Error("Failed to get user: %s error: %v", user_id, err)
			return nil, err
		}
		s.users[user_id] = user_obj
		return user_obj, nil
	} else {
		return user_obj, nil
	}
}

func (s *App) SaveUser(user *user.User) error {
	user_directory := path.Join(s.config.Database.UsersDir, user.Id)
	user_file_path := path.Join(user_directory, "user.json")

	err := ensureDirectoryExists(user_directory)
	if err != nil {
		s.logger.Error("Failed to create user directory: %s error: %v", user_directory, err)
		return err
	}

	file, err := os.Create(user_file_path)
	if err != nil {
		s.logger.Error("Failed to create user file: %s error: %v", user_file_path, err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(user)
	if err != nil {
		s.logger.Error("Failed save to user file: error: %v", err)
		return err
	}

	return nil
}

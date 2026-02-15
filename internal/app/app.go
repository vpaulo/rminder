package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

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

func (s *App) validateUserPath(user_id string) (string, error) {
	baseDir, err := filepath.Abs(s.config.Database.UsersDir)
	if err != nil {
		s.logger.Error("Failed to resolve base directory: %v", baseDir, err)
		return "", err
	}

	user_directory := filepath.Join(baseDir, user_id)
	resolved, err := filepath.Abs(user_directory)
	if err != nil {
		s.logger.Error("Failed to resolve user directory for user: %s error: %v", user_id, err)
		return "", err
	}

	if !strings.HasPrefix(resolved, baseDir+string(filepath.Separator)) {
		s.logger.Error("Directory traversal detected for user: %s resolved: %s", user_id, resolved)
		return "", err
	}

	return resolved, nil
}

func (s *App) GetDatabaseForUser(user_id string) (database.Service, error) {
	if db, ok := s.user_databases[user_id]; !ok {
		user_directory, err := s.validateUserPath(user_id)
		if err != nil {
			s.logger.Error("Invalid user path: %s error: %v", user_id, err)
			return nil, err
		}

		err = ensureDirectoryExists(user_directory)
		if err != nil {
			s.logger.Error("Failed to create directory: %s error: %v", user_directory, err)
			return nil, err
		}

		user_database_path := filepath.Join(user_directory, "db.sqlite")
		db := database.New(user_database_path)
		s.user_databases[user_id] = db
		return db, nil
	} else {
		return db, nil
	}
}

func (s *App) loadUserFromFile(user_id string) (*user.User, error) {
	user_directory, err := s.validateUserPath(user_id)
	if err != nil {
		s.logger.Error("Invalid user path: %s error: %v", user_id, err)
		return nil, err
	}
	user_file_path := filepath.Join(user_directory, "user.json")

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
	user_directory, err := s.validateUserPath(user.Id)
	if err != nil {
		s.logger.Error("Invalid user path: %s error: %v", user.Id, err)
		return err
	}
	user_file_path := filepath.Join(user_directory, "user.json")

	err = ensureDirectoryExists(user_directory)
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
		s.logger.Error("Failed save to user file", "error", err)
		return err
	}

	return nil
}

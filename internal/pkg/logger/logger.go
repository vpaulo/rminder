package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// Logger wraps slog.Logger
type Logger struct {
	*slog.Logger
	file *os.File
}

// Config holds logger configuration
type Config struct {
	Level  string
	Format string
	Output string
}

// New creates a new logger instance using slog
func New(cfg Config) (*Logger, error) {
	// Parse log level
	level := parseLevel(cfg.Level)

	// Configure output
	var output *os.File
	if cfg.Output != "" && cfg.Output != "stdout" {
		// Create parent directory if it does not exist
		dir := filepath.Dir(cfg.Output)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
		file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		output = file
	} else {
		output = os.Stdout
	}

	// Create handler based on format
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: level,
	}

	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	logger := slog.New(handler)

	return &Logger{Logger: logger, file: output}, nil
}

// parseLevel parses string level to slog.Level
func parseLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{Logger: l.Logger.With(key, value), file: l.file}
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{Logger: l.Logger.With(args...), file: l.file}
}

// WithUserID adds user_id to the logger context
func (l *Logger) WithUserID(userID string) *Logger {
	return l.WithField("user_id", userID)
}

// WithRequestID adds request_id to the logger context
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithField("request_id", requestID)
}

// Close syncs and closes the log file if one is open
func (l *Logger) Close() error {
	if l.file != nil && l.file != os.Stdout && l.file != os.Stderr {
		if err := l.file.Sync(); err != nil {
			return err
		}
		return l.file.Close()
	}
	return nil
}

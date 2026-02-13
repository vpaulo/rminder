package database

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

func (s *service) migrate() error {
	err := s.createMigrationsTable()
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	applied, err := s.appliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	migrations, err := availableMigrations()
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	for _, m := range migrations {
		if applied[m.version] {
			continue
		}

		log.Printf("applying migration %03d: %s", m.version, m.name)

		err = s.applyMigration(m)
		if err != nil {
			return fmt.Errorf("failed to apply migration %03d: %w", m.version, err)
		}
	}

	return nil
}

func (s *service) createMigrationsTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func (s *service) appliedMigrations() (map[int]bool, error) {
	rows, err := s.db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, rows.Err()
}

type migration struct {
	version  int
	name     string
	filename string
}

func availableMigrations() ([]migration, error) {
	entries, err := Migrations.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	var migrations []migration
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		// Parse version from filename: "001_initial_schema.sql" -> 1
		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) < 2 {
			continue
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		name := strings.TrimSuffix(parts[1], ".sql")
		migrations = append(migrations, migration{
			version:  version,
			name:     name,
			filename: entry.Name(),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].version < migrations[j].version
	})

	return migrations, nil
}

func (s *service) applyMigration(m migration) error {
	content, err := Migrations.ReadFile("migrations/" + m.filename)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(string(content))
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO schema_migrations (version, name) VALUES (?, ?)", m.version, m.name)
	if err != nil {
		return err
	}

	return tx.Commit()
}

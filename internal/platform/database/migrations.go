package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Run Table Migration
func RunMigrations(ctx context.Context, db *sql.DB, dir string) error {
	_, err := db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS schema_migrations (
            version TEXT PRIMARY KEY
        )`,
	)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		var exists bool
		err := db.QueryRowContext(
			ctx,
			"SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version=$1)",
			name,
		).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			continue
		}

		path := filepath.Join(dir, name)
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		fmt.Println("Running migration:", name)

		_, err = db.ExecContext(
			ctx,
			string(sqlBytes),
		)
		if err != nil {
			return err
		}

		_, err = db.Exec(
			"INSERT INTO schema_migrations (version) VALUES ($1)",
			name,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

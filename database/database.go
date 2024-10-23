package database

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

type migration struct {
	Name    string
	Hash    string
	Content []byte
}

func RunMigrations(db *sql.DB, dir string) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_history (
			name TEXT PRIMARY KEY NOT NULL,
			hash TEXT NOT NULL,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errors.New("no migration files found")
	}

	migrations := make(map[string]migration, len(files))
	for _, file := range files {
		if file.IsDir() {
			return errors.New("directory found in migrations folder, expected SQL files")
		}
		if filepath.Ext(file.Name()) != ".sql" {
			return errors.New("non-SQL file found in migrations folder")
		}

		content, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return err
		}

		f := migration{Name: file.Name(), Hash: fmt.Sprintf("%x", md5.Sum(content)), Content: content}
		migrations[f.Name] = f
	}

	rows, err := db.Query("SELECT name, hash FROM schema_history;")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name, hash string
		if err := rows.Scan(&name, &hash); err != nil {
			return err
		}

		migration, ok := migrations[name]
		if !ok {
			return fmt.Errorf("migration file %s missing", name)
		}
		if migration.Hash != hash {
			return fmt.Errorf("hash mismatch for migration %s", name)
		}
		delete(migrations, name)
	}

	if len(migrations) == 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO schema_history (name, hash) VALUES (?, ?);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, migration := range migrations {
		if _, err := tx.Exec(string(migration.Content)); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Name, err)
		}
		if _, err := stmt.Exec(migration.Name, migration.Hash); err != nil {
			return fmt.Errorf("failed to save migration %s in schema history: %w", migration.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

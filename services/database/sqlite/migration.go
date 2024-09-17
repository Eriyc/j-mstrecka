package sqlite

import (
	"database/sql"
	sqlite_migrations "gostrecka/services/database/sqlite/migrations"
	"strings"
)

func SetupMigrations(conn *sql.DB) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);
	`)

	if err != nil {
		return err
	}

	return nil
}

type Migration struct {
	Name    string
	Content string
}

func (m SqliteMiddleware) Migrate(conn *sql.DB) error {
	err := SetupMigrations(conn)
	if err != nil {
		m.Logger.Error("Error setting up migrations", "error", err.Error())
		return err
	}

	var files []Migration = []Migration{
		{Name: "20240809205925_initial", Content: sqlite_migrations.MIGRATION1},
	}

	for _, file := range files {

		row := conn.QueryRow("SELECT name FROM migrations WHERE name = ?", file.Name)
		var name string
		err = row.Scan(&name)

		if err == nil {
			m.Logger.Info("Skipping migration", "name", file.Name)
			continue
		}
		err = nil

		str := string(file.Content)
		stmts := strings.Split(str, ";\n\n")
		tx, err := conn.Begin()
		if err != nil {
			return err
		}

		for _, stmt := range stmts {
			_, err = tx.Exec(stmt)
			m.Logger.Info("Executing migration", "name", file.Name, "statement", stmt)
			if err != nil {
				tx.Rollback()
				m.Logger.Error("Error applying migration", "name", file.Name, "statement", stmt, "error", err.Error())
				return err
			}
		}

		_, err = tx.Exec("INSERT INTO migrations (name) VALUES (?)", file.Name)

		if err != nil {
			return err
		}

		tx.Commit()
	}
	return err
}

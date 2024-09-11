package sqlite

import (
	"database/sql"
	sqlite_migrations "gostrecka/internal/service/database/sqlite/migrations"
	"log"
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

func Migrate(conn *sql.DB) error {
	err := SetupMigrations(conn)
	if err != nil {
		log.Default().Fatalf("Error setting up migrations: %s", err)
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
			log.Printf("Skipping migration %s", file.Name)
			continue
		}
		str := string(file.Content)
		stmts := strings.Split(str, ";\n\n")
		tx, err := conn.Begin()
		if err != nil {
			return err
		}

		for _, stmt := range stmts {
			_, err = tx.Exec(stmt)
			log.Default().Printf("Applying migration %s: %s", file.Name, stmt)
			if err != nil {
				tx.Rollback()
				log.Fatalf("Error applying migration %s: %s", file.Name, stmt)
				return err
			}
		}

		_, err = tx.Exec("INSERT INTO migrations (name) VALUES (?)", file.Name)

		if err != nil {
			return err
		}

		tx.Commit()
	}

	println("Migrations applied")

	return err
}

package sqlite

import (
	"database/sql"
	"log"
	"os"
	"path"
	"slices"
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

func Migrate(conn *sql.DB) error {
	err := SetupMigrations(conn)
	if err != nil {
		log.Default().Fatalf("Error setting up migrations: %s", err)
		return err
	}

	var migrationPath = "./internal/service/database/sqlite/migrations"
	// get from folder
	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	log.Printf("Found %d migrations", len(files))

	commited, err := conn.Query("SELECT name FROM migrations")

	if err != nil {
		return err
	}

	var applied []string
	for commited.Next() {
		var name string
		err = commited.Scan(&name)
		if err != nil {
			return err
		}
		applied = append(applied, name)

	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		filePath := path.Join(migrationPath, fileName)

		fileName = strings.Split(fileName, ".")[0]
		if slices.Contains(applied, fileName) {
			continue
		}

		log.Printf("Applying migration %s", fileName)
		b, err := os.ReadFile(filePath)

		if err != nil {
			return err
		}
		str := string(b)
		stmts := strings.Split(str, ";\n\n")
		tx, err := conn.Begin()
		if err != nil {
			return err
		}

		for _, stmt := range stmts {
			_, err = tx.Exec(stmt)
			if err != nil {
				tx.Rollback()
				log.Fatalf("Error applying migration %s: %s", fileName, stmt)
				return err
			}
		}

		_, err = tx.Exec("INSERT INTO migrations (name) VALUES (?)", fileName)

		if err != nil {
			return err
		}

		err = tx.Commit()

		return err

	}

	println("Migrations applied")

	return err
}

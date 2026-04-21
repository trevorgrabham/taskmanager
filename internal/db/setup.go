package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"strconv"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func createMigrationTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
		);`)

	return err
}

func applyMigration(db *sql.DB, version int, fileName string) error {
	var (
		tx                                         *sql.Tx
		err                                        error
		versionCount, countAppliedFutureMigrations int
		migrationOperation                         []byte
	)
	if tx, err = startDBSession(db); err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = ?`, version).Scan(&versionCount)
	if err != nil {
		return fmt.Errorf("Migrate: %s", err)
	}

	if versionCount > 0 {
		return nil
	}

	// unapplied migration
	err = tx.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version > ?`, version).Scan(&countAppliedFutureMigrations)
	if err != nil {
		return fmt.Errorf("Migrate: %s", err)
	}
	if countAppliedFutureMigrations > 0 {
		log.Fatalf("migration state corrupted: version %d unapplied, but future version is applied", version)
	}

	migrationOperation, err = migrationFiles.ReadFile(fmt.Sprintf("migrations/%s", fileName))
	if err != nil {
		return fmt.Errorf("Migrate: %s", err)
	}

	if _, err = tx.Exec(string(migrationOperation)); err != nil {
		return fmt.Errorf("Migrate: %s", err)
	}

	if _, err = tx.Exec(`INSERT INTO schema_migrations(version) VALUES (?)`, version); err != nil {
		return fmt.Errorf("Migrate: %s", err)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func Migrate(db *sql.DB) error {
	var (
		err           error
		dirEntries    []fs.DirEntry
		versionString string
		version       int
	)
	err = createMigrationTable(db)
	if err != nil {
		return fmt.Errorf("Migrate: %s", err)
	}

	dirEntries, err = migrationFiles.ReadDir("migrations") // returns dirEntries in sorted (alphabetical) order. So versions appear in the order they should be applied
	if err != nil {
		return fmt.Errorf("Migrate: %s", err)
	}

	for _, entry := range dirEntries {
		versionString = entry.Name()[:3] // of the form "XXX-description.sql". i.e. "001-initial-tasks-and-recurring-tables.sql"
		version, err = strconv.Atoi(versionString)
		if err != nil {
			return fmt.Errorf("Migrate: %s", err)
		}

		if err = applyMigration(db, version, entry.Name()); err != nil {
			return fmt.Errorf("Migrate: %s", err)
		}

	}
	return nil
}

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, fmt.Errorf("connecting: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("connecting: %s", err)
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		return nil, fmt.Errorf("connecting: %s", err)
	}

	if err = Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// createMigrationTableIfNotExists creates a table to track applied migrations if it doesn't already exist.
func createMigrationTableIfNotExists(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
		);`)

	return err
}

// applyMigration applies the SQL commands in the .sql file identified by fileName if it has not already been applied.
// Propogates errors back to the calling function.
//
// Exits if the migration state becomes corrupted. This would happen by migrations being applied in the wrong order, such as having a future migration applied before the current one under inspection.
func applyMigration(db *sql.DB, fileName string) error {
	var (
		tx                                                  *sql.Tx
		err                                                 error
		version, versionCount, countAppliedFutureMigrations int
		migrationOperation                                  []byte
	)
	// of the form "XXX-description.sql". i.e. "001-initial-tasks-and-recurring-tables.sql"
	version, err = strconv.Atoi(fileName[:3])
	if err != nil {
		return err
	}

	if tx, err = db.Begin(); err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = ?`, version).Scan(&versionCount)
	if err != nil {
		return err
	}

	if versionCount > 0 {
		return nil
	}

	// unapplied migration
	err = tx.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version > ?`, version).Scan(&countAppliedFutureMigrations)
	if err != nil {
		return err
	}
	if countAppliedFutureMigrations > 0 {
		log.Fatalf("migration state corrupted: version %d unapplied, but future version is applied", version)
	}

	migrationOperation, err = migrationFiles.ReadFile(fmt.Sprintf("migrations/%s", fileName))
	if err != nil {
		return err
	}

	if _, err = tx.Exec(string(migrationOperation)); err != nil {
		return err
	}

	if _, err = tx.Exec(`INSERT INTO schema_migrations(version) VALUES (?)`, version); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Migrate applies any migration files that have not yet been applied.
func Migrate(db *sql.DB) (err error) {
	var dirEntries []fs.DirEntry
	err = createMigrationTableIfNotExists(db)
	if err != nil {
		return fmt.Errorf("Migrate: %w", err)
	}

	dirEntries, err = migrationFiles.ReadDir("migrations") // returns dirEntries in sorted (alphabetical) order. So versions appear in the order they should be applied
	if err != nil {
		return fmt.Errorf("Migrate: %s", err)
	}

	for _, entry := range dirEntries {
		if err = applyMigration(db, entry.Name()); err != nil {
			return fmt.Errorf("Migrate: %s", err)
		}

	}
	return nil
}

// Connect opens a database connection pool for the Repo and applies any unapplied migrations.
func (r *Repo) Connect(file string) error {
	var (
		caller = "Connect"
		err    error
	)

	r.db, err = sql.Open("sqlite3", file)
	if err != nil {
		return fmt.Errorf("%s: %w: %s", caller, ErrInternalRepo, err)
	}

	err = r.db.Ping()
	if err != nil {
		return fmt.Errorf("%s: %w: %s", caller, ErrInternalRepo, err)
	}

	_, err = r.db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		return fmt.Errorf("%s: %w: %s", caller, ErrInternalRepo, err)
	}

	if err = Migrate(r.db); err != nil {
		return fmt.Errorf("%s: %w: %s", caller, ErrInternalRepo, err)
	}

	return nil
}

func (r *Repo) Close() {
	r.db.Close()
}

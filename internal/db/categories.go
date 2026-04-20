package db

import (
	"database/sql"
)

func Categories(db *sql.DB) ([]string, error) {
	if err := checkDBConnection(db); err != nil { return nil, err }

	return getCategories(db)
}

package db

import (
	"database/sql"
	"fmt"
)

func Categories(db *sql.DB) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("categories: cannot get categories for a nil database")
	}

	rows, err := db.Query(`SELECT DISTINCT category FROM task ORDER BY category`)
	if err != nil {
		return nil, fmt.Errorf("catgories: %s", err)
	}
	defer rows.Close()

	var (
		category   string
		categories []string
	)
	for rows.Next() {
		err = rows.Scan(&category)
		if err != nil {
			return nil, fmt.Errorf("catgories: %s", err)
		}

		categories = append(categories, category)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("catgories: %s", err)
	}

	return categories, nil
}

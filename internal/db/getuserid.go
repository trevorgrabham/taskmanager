package db

import (
	"database/sql"
	"fmt"
)

func GetUserID(db *sql.DB, sessionID string) (userID int, err error) {
	if err = checkDBConnection(db); err != nil { return 0, err }
	if sessionID == "" { return 0, fmt.Errorf("GetUserID: no sessionID") }

	return getUserID(db, sessionID)
}

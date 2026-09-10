package db

import (
	"database/sql"

	"school-app/database"
)

func Open(filename string) (*sql.DB, error) {
	return database.Open(filename)
}

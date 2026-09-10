package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Open(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return nil, err
	}

	if _, err := db.Exec(`PRAGMA journal_mode = WAL`); err != nil {
		return nil, err
	}

	if _, err := db.Exec(`PRAGMA busy_timeout = 5000`); err != nil {
		return nil, err
	}

	if err := createSchema(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS students (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			nickname TEXT,
			parent_name TEXT,
			address1 TEXT,
			address2 TEXT,
			city TEXT,
			state TEXT,
			postal_code TEXT,
			email TEXT,
			grade TEXT
		);
	`)

	return err
}

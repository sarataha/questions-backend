package util

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
        -- content of migrations/init.sql
    `)
	if err != nil {
		return nil, err
	}
	return db, nil
}

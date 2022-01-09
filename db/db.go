package db

import (
	"database/sql"
	"errors"
)

var db *sql.DB

func InitDb() error {
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/nhk?parseTime=true")

	if err != nil {
		return err
	}

	return nil
}

func CloseDb() error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.Close()
}

func Query(query string) (*sql.Rows, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	return db.Query(query)
}

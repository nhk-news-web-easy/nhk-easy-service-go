package db

import (
	"database/sql"
	"errors"
	"fmt"
)

var db *sql.DB

func InitDb(config DbConfig) error {
	var err error
	db, err = sql.Open(config.DriverName, fmt.Sprintf("%s:%s@tcp(%s:%d)/nhk?parseTime=true", config.UserName, config.Password, config.Host, config.Port))

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

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	return db.Query(query, args...)
}

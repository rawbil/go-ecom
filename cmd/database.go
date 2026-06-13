package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func dbConnection(config DBConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("Database Connected")

	return db, nil
}

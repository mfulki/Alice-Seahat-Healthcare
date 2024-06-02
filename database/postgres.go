package database

import (
	"database/sql"
	"fmt"

	"Alice-Seahat-Healthcare/seahat-be/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnPostgres() (*sql.DB, error) {
	dbURL := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.DB.Host,
		config.DB.User,
		config.DB.Password,
		config.DB.Name,
		config.DB.Port,
	)

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

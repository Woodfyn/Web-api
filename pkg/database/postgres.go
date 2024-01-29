package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ConnInfo struct {
	Host     string
	Port     string
	Username string
	Name     string
	SSLMode  string
	Password string
}

func NewPostgesDB(cfg ConnInfo) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Name, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

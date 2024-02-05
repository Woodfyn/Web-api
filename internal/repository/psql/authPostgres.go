package psql

import (
	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UsersPostgres struct {
	db *sqlx.DB
}

func NewUsers(db *sqlx.DB) *UsersPostgres {
	return &UsersPostgres{
		db: db,
	}
}

func (r *UsersPostgres) CreateUser(user domain.SignUpInput) (int, error) {
	return 0, nil
}

func (r *UsersPostgres) GetUser(username, password string) (domain.User, error) {
	return domain.User{}, nil
}

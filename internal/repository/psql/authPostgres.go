package psql

import (
	"fmt"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

const userTable = "users"

type UsersPostgres struct {
	db *sqlx.DB
}

func NewUsers(db *sqlx.DB) *UsersPostgres {
	return &UsersPostgres{
		db: db,
	}
}

func (r *UsersPostgres) Create(inp domain.User) error {
	_, err := r.db.Exec(fmt.Sprintf("INSERT INTO %s (name, email, password, registered_at) VALUES ($1, $2, $3, $4)", userTable),
		inp.Name, inp.Email, inp.Password, inp.RegisteredAt)

	return err
}

func (r *UsersPostgres) GetByCredentials(email, password string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(fmt.Sprintf("SELECT id, name, email, registered_at FROM %s WHERE email=$1 AND password=$2", userTable),
		email, password).
		Scan(&user.ID, &user.Name, &user.Email, &user.RegisteredAt)

	return user, err
}

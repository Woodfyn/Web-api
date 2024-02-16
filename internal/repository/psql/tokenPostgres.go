package psql

import (
	"fmt"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

const tokenTable = "refresh_tokens"

type Tokens struct {
	db *sqlx.DB
}

func NewTokens(db *sqlx.DB) *Tokens {
	return &Tokens{db}
}

func (r *Tokens) Create(token domain.RefreshSession) error {
	_, err := r.db.Exec(fmt.Sprintf("INSERT INTO %s (user_id, token, expires_at) values ($1, $2, $3)", tokenTable),
		token.UserID, token.Token, token.ExpiresAt)

	return err
}

func (r *Tokens) Get(token string) (domain.RefreshSession, error) {
	var t domain.RefreshSession
	err := r.db.QueryRow(fmt.Sprintf("SELECT id, user_id, token, expires_at FROM %s WHERE token=$1", tokenTable), token).
		Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt)
	if err != nil {
		return t, err
	}

	_, err = r.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", tokenTable), t.UserID)

	return t, err
}

func (r *Tokens) GetByRefreshToken(token string) (domain.RefreshSession, error) {
	var t domain.RefreshSession
	err := r.db.QueryRow(fmt.Sprintf("SELECT id, user_id, token, expires_at FROM %s WHERE token=$1", tokenTable), token).
		Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt)
	if err != nil {
		return t, err
	}

	return t, err
}

func (r *Tokens) Delete(token string) error {
	_, err := r.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE token = $1", tokenTable), token)
	return err
}

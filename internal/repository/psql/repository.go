package psql

import (
	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TokenSessions interface {
	Create(token domain.RefreshSession) error
	Get(token string) (domain.RefreshSession, error)
	GetByRefreshToken(token string) (domain.RefreshSession, error)
	Delete(token string) error
}

type Users interface {
	Create(inp domain.User) error
	Get(id int) (domain.User, error)
	GetByCredentials(email, password string) (domain.User, error)
}

type Games interface {
	Create(game domain.Game) error
	GetAll() ([]domain.Game, error)
	GetById(gameId int) (domain.Game, error)
	UpdateById(gameId int, input domain.UpdateGameInput) error
	DeleteById(gameId int) error
}

type Repositories struct {
	Users  Users
	Tokens TokenSessions
	Games  Games
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users:  NewUsers(db),
		Tokens: NewTokens(db),
		Games:  NewGames(db),
	}
}

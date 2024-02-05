package psql

import (
	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Users interface {
	CreateUser(user domain.SignUpInput) (int, error)
	GetUser(username, password string) (domain.User, error)
}

type Games interface {
	Create(game domain.Game) error
	GetAll() ([]domain.Game, error)
	GetById(gameId int) (domain.Game, error)
	UpdateById(gameId int, input domain.UpdateGameInput) error
	DeleteById(gameId int) error
}

type Repositories struct {
	Users Users
	Games Games
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users: NewUsers(db),
		Games: NewGames(db),
	}
}

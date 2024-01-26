package psql

import (
	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Game interface {
	Create(game domain.Game) (int, error)
	GetAll() ([]domain.Game, error)
	GetById(gameId int) (domain.Game, error)
	UpdateById(gameId int, input domain.UpdateItemInput) error
	DeleteById(gameId int) error
}

type Repository struct {
	Game
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Game: NewGamePostgres(db),
	}
}

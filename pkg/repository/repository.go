package repository

import (
	todo "github.com/Woodfyn/Web-api"

	"github.com/jmoiron/sqlx"
)

type Game interface {
	Create(game todo.Game) (int, error)
	GetAll() ([]todo.Game, error)
	GetById(gameId int) (todo.Game, error)
	UpdateById(gameId int, input todo.UpdateItemInput) error
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

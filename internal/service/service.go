package service

import (
	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/Woodfyn/Web-api/internal/repository/psql"
)

type Game interface {
	Create(game domain.Game) (int, error)
	GetAll() ([]domain.Game, error)
	GetById(gameId int) (domain.Game, error)
	UpdateById(gameId int, input domain.UpdateItemInput) error
	DeleteById(gameId int) error
}

type Service struct {
	Game
}

func NewService(repos *psql.Repository) *Service {
	return &Service{
		Game: NewGameService(repos.Game),
	}
}

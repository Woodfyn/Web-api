package service

import (
	todo "github.com/Woodfyn/Web-api"

	"github.com/Woodfyn/Web-api/pkg/repository"
)

type Game interface {
	Create(game todo.Game) (int, error)
	GetAll() ([]todo.Game, error)
	GetById(gameId int) (todo.Game, error)
	UpdateById(gameId int, input todo.UpdateItemInput) error
	DeleteById(gameId int) error
}

type Service struct {
	Game
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Game: NewGameService(repos.Game),
	}
}

package service

import (
	"github.com/Woodfyn/Web-api/pkg/repository"

	todo "github.com/Woodfyn/Web-api"
)

type GameService struct {
	repo repository.Game
}

func NewGameService(repo repository.Game) *GameService {
	return &GameService{repo: repo}
}

func (s *GameService) Create(game todo.Game) (int, error) {
	return s.repo.Create(game)
}

func (s *GameService) GetAll() ([]todo.Game, error) {
	return s.repo.GetAll()
}

func (s *GameService) GetById(gameId int) (todo.Game, error) {
	return s.repo.GetById(gameId)
}

func (s *GameService) UpdateById(gameId int, input todo.UpdateItemInput) error {
	return s.repo.UpdateById(gameId, input)
}

func (s *GameService) DeleteById(gameId int) error {
	return s.repo.DeleteById(gameId)
}

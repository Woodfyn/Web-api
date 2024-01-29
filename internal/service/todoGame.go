package service

import (
	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/Woodfyn/Web-api/internal/repository/psql"
)

type GameService struct {
	repo psql.Game
}

func NewGameService(repo psql.Game) *GameService {
	return &GameService{repo: repo}
}

func (s *GameService) Create(game domain.Game) (int, error) {
	return s.repo.Create(game)
}

func (s *GameService) GetAll() ([]domain.Game, error) {
	return s.repo.GetAll()
}

func (s *GameService) GetById(gameId int) (domain.Game, error) {
	return s.repo.GetById(gameId)
}

func (s *GameService) UpdateById(gameId int, input domain.UpdateItemInput) error {
	return s.repo.UpdateById(gameId, input)
}

func (s *GameService) DeleteById(gameId int) error {
	return s.repo.DeleteById(gameId)
}

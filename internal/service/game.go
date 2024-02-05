package service

import (
	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/Woodfyn/Web-api/internal/repository/psql"
)

type Game struct {
	repo psql.Games
}

func NewServiceGame(repo psql.Games) *Game {
	return &Game{
		repo: repo,
	}
}

func (s *Game) Create(game domain.Game) error {
	return s.repo.Create(game)
}

func (s *Game) GetAll() ([]domain.Game, error) {
	return s.repo.GetAll()
}

func (s *Game) GetById(gameId int) (domain.Game, error) {
	return s.repo.GetById(gameId)
}

func (s *Game) UpdateById(gameId int, input domain.UpdateGameInput) error {
	return s.repo.UpdateById(gameId, input)
}

func (s *Game) DeleteById(gameId int) error {
	return s.repo.DeleteById(gameId)
}

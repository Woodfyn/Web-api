package service

import (
	"context"
	"time"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/Woodfyn/Web-api/internal/repository/psql"
	"github.com/Woodfyn/auditLog/pkg/core"
	"github.com/sirupsen/logrus"
)

type Game struct {
	repo  psql.Games
	audit AuditClient
}

func NewServiceGame(repo psql.Games, audit AuditClient) *Game {
	return &Game{
		repo:  repo,
		audit: audit,
	}
}

func (s *Game) Create(game domain.Game) error {
	if err := s.audit.SendLogRequest(context.TODO(), core.LogItem{
		Action:    core.ACTION_CREATE,
		Entity:    core.ENTITY_GAME,
		EntityID:  int64(game.Id),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "Create.Audit",
		}).Error("failed to send log request:", err)
	}

	return s.repo.Create(game)
}

func (s *Game) GetAll() ([]domain.Game, error) {
	if err := s.audit.SendLogRequest(context.TODO(), core.LogItem{
		Action:    core.ACTION_GET,
		Entity:    core.ENTITY_GAME,
		EntityID:  0,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "GetAll.Audit",
		}).Error("failed to send log request:", err)
	}

	return s.repo.GetAll()
}

func (s *Game) GetById(gameId int) (domain.Game, error) {
	if err := s.audit.SendLogRequest(context.TODO(), core.LogItem{
		Action:    core.ACTION_GET,
		Entity:    core.ENTITY_GAME,
		EntityID:  int64(gameId),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "GetById.Audit",
		}).Error("failed to send log request:", err)
	}

	return s.repo.GetById(gameId)
}

func (s *Game) Update(gameId int, input domain.UpdateGameInput) error {
	if err := s.audit.SendLogRequest(context.TODO(), core.LogItem{
		Action:    core.ACTION_UPDATE,
		Entity:    core.ENTITY_GAME,
		EntityID:  int64(gameId),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "Update.Audit",
		}).Error("failed to send log request:", err)
	}

	return s.repo.UpdateById(gameId, input)
}

func (s *Game) Delete(gameId int) error {
	if err := s.audit.SendLogRequest(context.TODO(), core.LogItem{
		Action:    core.ACTION_DELETE,
		Entity:    core.ENTITY_GAME,
		EntityID:  int64(gameId),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "Delete.Audit",
		}).Error("failed to send log request:", err)
	}

	return s.repo.DeleteById(gameId)
}

package service

import (
	"context"
	"strconv"
	"time"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/Woodfyn/Web-api/internal/repository/psql"
	"github.com/Woodfyn/Web-api/pkg/auth"
	audit "github.com/Woodfyn/auditLog/pkg/core"
	"github.com/sirupsen/logrus"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type User struct {
	repo        psql.Users
	sessionRepo psql.TokenSessions
	hasher      PasswordHasher

	audit           AuditClient
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewServiceUser(repo psql.Users, sessionRepo psql.TokenSessions, hasher PasswordHasher, audit AuditClient, tokenManager auth.TokenManager, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *User {
	return &User{
		repo:        repo,
		sessionRepo: sessionRepo,
		hasher:      hasher,

		audit:           audit,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *User) SignUp(inp domain.SignUpInput) error {
	hashPassword, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     hashPassword,
		RegisteredAt: time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return err
	}

	user, err = s.repo.GetByCredentials(inp.Email, hashPassword)
	if err != nil {
		return err
	}

	if err := s.audit.SendLogRequest(context.TODO(), audit.LogItem{
		Action:    audit.ACTION_REGISTER,
		Entity:    audit.ENTITY_USER,
		EntityID:  int64(user.ID),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "SignUp.Audit",
		}).Error("failed to send log request:", err)
	}

	return nil
}

func (s *User) SignIn(inp domain.SignInInput) (string, string, error) {
	hashPassword, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetByCredentials(inp.Email, hashPassword)
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := s.genereteTokens(user.ID)
	if err != nil {
		return "", "", err
	}

	if err := s.audit.SendLogRequest(context.TODO(), audit.LogItem{
		Action:    audit.ACTION_LOGIN,
		Entity:    audit.ENTITY_USER,
		EntityID:  int64(user.ID),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "SignIn.Audit",
		}).Error("failed to send log request:", err)
	}

	return accessToken, refreshToken, nil
}

func (s *User) RefreshTokens(refreshToken string) (string, string, error) {
	session, err := s.sessionRepo.Get(refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", domain.ErrRefreshTokenExpired
	}

	if err := s.audit.SendLogRequest(context.TODO(), audit.LogItem{
		Action:    audit.ACTION_LOGIN,
		Entity:    audit.ENTITY_USER,
		EntityID:  int64(session.ID),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "RefreshTokens.Audit",
		}).Error("failed to send log request:", err)
	}

	return s.genereteTokens(session.UserID)
}

func (s *User) genereteTokens(userId int) (string, string, error) {
	accessToken, err := s.tokenManager.NewJWT(strconv.Itoa(userId), s.accessTokenTTL)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.sessionRepo.Create(domain.RefreshSession{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *User) LogOut(refreshToken string) error {
	refreshSession, err := s.sessionRepo.GetByRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	user, err := s.repo.Get(refreshSession.UserID)
	if err != nil {
		return err
	}

	if err := s.audit.SendLogRequest(context.TODO(), audit.LogItem{
		Action:    audit.ACTION_LOGOUT,
		Entity:    audit.ENTITY_USER,
		EntityID:  int64(user.ID),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "SignUp.Audit",
		}).Error("failed to send log request:", err)
	}

	return s.sessionRepo.Delete(refreshToken)
}

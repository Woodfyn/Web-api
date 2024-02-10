package service

import (
	"time"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/Woodfyn/Web-api/internal/repository/psql"
	"github.com/gorilla/sessions"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type User struct {
	repo   psql.Users
	hasher PasswordHasher
}

func NewServiceUser(repo psql.Users, hasher PasswordHasher) *User {
	return &User{
		repo:   repo,
		hasher: hasher,
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

	return nil
}

func (s *User) SignIn(inp domain.SignInInput, session *sessions.Session) (*sessions.Session, error) {
	hashPassword, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetByCredentials(inp.Email, hashPassword)
	if err != nil {
		return nil, err
	}

	session.Values["user_id"] = user.ID

	return session, err
}

func (s *User) LogOut(session *sessions.Session) (*sessions.Session, error) {
	session.Options.MaxAge = -1

	return session, nil
}

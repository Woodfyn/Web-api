package service

import (
	"time"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/Woodfyn/Web-api/internal/repository/psql"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserService struct {
	repo   psql.UsersPostgres
	hasher PasswordHasher

	hmacSecret []byte
	tokenTtl   time.Duration
}

func NewServiceUser(repo psql.UsersPostgres, hasher PasswordHasher, hmacSecret []byte, tokenTtl time.Duration) *UserService {
	return &UserService{
		repo:       repo,
		hasher:     hasher,
		hmacSecret: hmacSecret,
		tokenTtl:   tokenTtl,
	}
}

func (s *UserService) CreateUser(user domain.SignUpInput) error {
	return nil
}

func (s *UserService) GenerateToken(username, password string) (string, error) {
	return "", nil
}

func (s *UserService) ParseToken(token string) (int, error) {
	return 0, nil
}

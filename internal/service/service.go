package service

import (
	"time"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/Woodfyn/Web-api/internal/repository/psql"
	"github.com/Woodfyn/Web-api/pkg/auth"
	"github.com/Woodfyn/Web-api/pkg/hash"
	"github.com/Woodfyn/auditLog/pkg/core"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Games interface {
	Create(game domain.Game) error
	GetAll() ([]domain.Game, error)
	GetById(gameId int) (domain.Game, error)
	Update(gameId int, input domain.UpdateGameInput) error
	Delete(gameId int) error
}

type Users interface {
	SignUp(user domain.SignUpInput) error
	SignIn(inp domain.SignInInput) (string, string, error)
	RefreshTokens(refreshToken string) (string, string, error)
	LogOut(refreshToken string) error
}

type MQClient interface {
	Publisher(item core.LogItem) error
}

type Services struct {
	Games Games
	Users Users
}

type Deps struct {
	Repos  *psql.Repositories
	Hasher hash.PasswordHasher

	MQClient        MQClient
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(deps Deps) *Services {
	return &Services{
		Games: NewServiceGame(deps.Repos.Games, deps.MQClient),
		Users: NewServiceUser(deps.Repos.Users, deps.Repos.Tokens, deps.Hasher, deps.MQClient, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
	}
}

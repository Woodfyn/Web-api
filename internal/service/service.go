package service

import (
	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/Woodfyn/Web-api/internal/repository/psql"
	"github.com/Woodfyn/Web-api/pkg/hash"
	"github.com/gorilla/sessions"
)

type Games interface {
	Create(game domain.Game) error
	GetAll() ([]domain.Game, error)
	GetById(gameId int) (domain.Game, error)
	Update(gameId int, input domain.UpdateGameInput) error
	Delete(gameId int) error
}

type Users interface {
	SignUp(user domain.SignUpInput) error
	SignIn(inp domain.SignInInput, session *sessions.Session) (*sessions.Session, error)
	LogOut(session *sessions.Session) (*sessions.Session, error)
}

type Services struct {
	Games Games
	Users Users
}

type Deps struct {
	Repos  *psql.Repositories
	Hasher hash.PasswordHasher
}

func NewServices(deps Deps) *Services {
	return &Services{
		Games: NewServiceGame(deps.Repos.Games),
		Users: NewServiceUser(deps.Repos.Users, deps.Hasher),
	}
}

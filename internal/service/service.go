package service

import (
	"time"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/Woodfyn/Web-api/internal/repository/psql"
	"github.com/Woodfyn/Web-api/pkg/hash"
)

type Games interface {
	Create(book domain.Game) error
	GetByID(id int) (domain.Game, error)
	GetAll() ([]domain.Game, error)
	Delete(id int) error
	Update(id int, inp domain.UpdateGameInput) error
}

type Users interface {
	Create(user domain.User) error
	GetByCredentials(email, password string) (domain.User, error)
}

type Services struct {
	Games Games
	Users Users
}

type Deps struct {
	Repos  *psql.Repositories
	Hasher hash.PasswordHasher

	hmacSecret []byte
	tokenTtl   time.Duration
}

func NewServices(deps Deps) *Services {
	return &Services{
		games: NewServiceGame(deps.Repos),
		users: NewServiceUser(deps.Repos, deps.Hasher, deps.hmacSecret, deps.tokenTtl),
	}
}

package domain

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

var ErrUserNotFound = errors.New("user with such credentials not found")

type User struct {
	ID           int
	Name         string
	Email        string
	Password     string
	RegisteredAt time.Time
}

type SignUpInput struct {
	Name     string `validate:"required,get=2"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,get=6"`
}

func (v *SignUpInput) Validate() error {
	return validate.Struct(v)
}

type SignInInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,get=6"`
}

func (v *SignInInput) Validate() error {
	return validate.Struct(v)
}

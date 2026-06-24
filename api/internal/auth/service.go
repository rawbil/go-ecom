package auth

import (
	"context"
	"database/sql"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
)

type Service interface {
	UserRegister(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error)

	// UserLogin()
	// UserLogout()
}

type Svc struct {
	repository repository.Queries
}

func NewService(repository repository.Queries) Service {
	return &Svc{
		repository: repository,
	}
}

func (svc *Svc) UserRegister(ctx context.Context, params repository.CreateUserParams) (sql.Result, error) {
	
	return svc.repository.CreateUser(ctx, params)
}

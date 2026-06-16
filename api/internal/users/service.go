package users

import (
	"context"
	"database/sql"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
)

type Service interface {
	ListAllUsers(ctx context.Context) (users []repository.User, err error)
	ListUser(ctx context.Context, email string) (user repository.User, err error)
	CreateUser(ctx context.Context, params repository.CreateUserParams) (result sql.Result, err error)
	DeleteUser(ctx context.Context, email string) error
}

type Svc struct {
	repository repository.Queries
}

func NewService(repository repository.Queries) Service {
	return &Svc{
		repository: repository,
	}
}

func (svc *Svc) ListAllUsers(ctx context.Context) (users []repository.User, err error) {
	return svc.repository.ListUsers(ctx)
}

func (svc *Svc) ListUser(ctx context.Context, email string) (user repository.User, err error) {
	return svc.repository.ListUser(ctx, email)
}

func (svc *Svc) CreateUser(ctx context.Context, params repository.CreateUserParams) (result sql.Result, err error) {
	return svc.repository.CreateUser(ctx, params)
}

func (svc *Svc) DeleteUser(ctx context.Context, email string) error {
	return svc.repository.DeleteUser(ctx, email)
}

package users

import (
	"context"
	"database/sql"
	"errors"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
)

type Service interface {
	ListAllUsers(ctx context.Context, arg repository.ListUsersParams) (users []repository.User, err error)
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

var (
	AuthNotFound = errors.New("User not found in Context. Try logging in again...")
	UserNotFound = errors.New("User Not Found. Try logging in again...")
)

func (svc *Svc) ListAllUsers(ctx context.Context, arg repository.ListUsersParams) (users []repository.User, err error) {
	return svc.repository.ListUsers(ctx, arg)
}

func (svc *Svc) ListUser(ctx context.Context, email string) (user repository.User, err error) {
	return svc.repository.ListUser(ctx, email)
}

func (svc *Svc) CreateUser(ctx context.Context, params repository.CreateUserParams) (result sql.Result, err error) {
	return svc.repository.CreateUser(ctx, params)
}

func (svc *Svc) DeleteUser(ctx context.Context, email string) error {
	user_id, ok := authutils.GetUserIDFromContext(ctx)
	if !ok {
		return AuthNotFound
	}

	//& Ensure user exists
	if _, err := svc.repository.ListUserById(ctx, user_id); err != nil {
		return UserNotFound
	}
	
	return svc.repository.DeleteUser(ctx, email)
}

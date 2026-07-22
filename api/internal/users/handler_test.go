package users

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	"github.com/rawbil/ecom2/internal/utils"
)

type MockUserService struct {
	ListAllUsersFunc func(ctx context.Context, arg repository.ListUsersParams) (users []repository.User, err error)
	ListUserFunc     func(ctx context.Context, email string) (user repository.User, err error)
	CreateUserFunc   func(ctx context.Context, params repository.CreateUserParams) (result sql.Result, err error)
	DeleteUserFunc   func(ctx context.Context, email string) error
}

func (m *MockUserService) ListAllUsers(ctx context.Context, arg repository.ListUsersParams) (users []repository.User, err error) {
	return m.ListAllUsersFunc(ctx, arg)
}

func (m *MockUserService) ListUser(ctx context.Context, email string) (user repository.User, err error) {
	return m.ListUserFunc(ctx, email)
}

func (m *MockUserService) CreateUser(ctx context.Context, params repository.CreateUserParams) (result sql.Result, err error) {
	return m.CreateUserFunc(ctx, params)
}

func (m *MockUserService) DeleteUser(ctx context.Context, email string) error {
	return m.DeleteUserFunc(ctx, email)
}

func TestListAllUsers(t *testing.T) {
	utils.Slogger()
	tests := []struct {
		name           string
		target         string
		serviceErr     error
		expectedStatus int
	}{
		{
			name:           "success",
			target:         "/users?email=bildadsimiyu6@gmail.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid email",
			target:         "/users?email=bildadsimiyu6",
			serviceErr:     InvalidEmail,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := &MockUserService{
				ListAllUsersFunc: func(ctx context.Context, arg repository.ListUsersParams) (users []repository.User, err error) {
					// if test.serviceErr != nil {
					// 	return []repository.User{}, test.serviceErr
					// }

					// if arg.Email == "" {
					// 	return []repository.User{}, errors.New("email query was not passed to service")
					// }

					// return a slice with one mocked user
					return []repository.User{{}}, test.serviceErr
				},
			}

			h := NewHandler(mock)

			req := httptest.NewRequest(http.MethodGet, test.target, nil)
			rr := httptest.NewRecorder()
			h.ListAllUsers(rr, req)

			if rr.Code != test.expectedStatus {
				t.Fatalf("expected %d, got %d", test.expectedStatus, rr.Code)
			}
		})
	}

}

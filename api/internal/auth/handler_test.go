package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
	"github.com/rawbil/ecom2/internal/utils"
)

type MockAuthService struct {
	UserRegisterFunc func(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error)

	UserLoginFunc        func(ctx context.Context, arg authutils.UserLoginParams) (repository.User, string, string, error)
	LogoutUserFunc       func(ctx context.Context) error
	PasswordResetFunc    func(ctx context.Context, arg authutils.PasswordResetParams) error
	RefreshTokensFunc    func(ctx context.Context, arg authutils.RefreshTokenParam) (string, string, error)
	UpdateMyEmailFunc    func(ctx context.Context, arg authutils.UpdateEmailParams) (repository.User, error)
	UpdateMyUsernameFunc func(ctx context.Context, arg authutils.UpdateUsernameParams) (repository.User, error)
	ForgotPasswordFunc   func(ctx context.Context, arg authutils.UpdateEmailParams) error
}

// * MOCK SERVICES
func (m *MockAuthService) UserRegister(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error) {
	return m.UserRegisterFunc(ctx, arg)
}

func (m *MockAuthService) UserLogin(ctx context.Context, arg authutils.UserLoginParams) (repository.User, string, string, error) {
	return m.UserLoginFunc(ctx, arg)
}

func (m *MockAuthService) LogoutUser(ctx context.Context) error {
	return m.LogoutUserFunc(ctx)
}

func (m *MockAuthService) PasswordReset(ctx context.Context, arg authutils.PasswordResetParams) error {
	return m.PasswordResetFunc(ctx, arg)
}

func (m *MockAuthService) RefreshTokens(ctx context.Context, arg authutils.RefreshTokenParam) (string, string, error) {
	return m.RefreshTokensFunc(ctx, arg)
}

func (m *MockAuthService) UpdateMyEmail(ctx context.Context, arg authutils.UpdateEmailParams) (repository.User, error) {
	return m.UpdateMyEmailFunc(ctx, arg)
}

func (m *MockAuthService) UpdateMyUsername(ctx context.Context, arg authutils.UpdateUsernameParams) (repository.User, error) {
	return m.UpdateMyUsernameFunc(ctx, arg)
}

func (m *MockAuthService) ForgotPassword(ctx context.Context, arg authutils.UpdateEmailParams) error {
	return m.ForgotPasswordFunc(ctx, arg)
}

// * Register Test
func TestRegisterHandler(t *testing.T) {
	utils.Slogger()

	tests := []struct {
		name            string
		expectedStatus  int
		responseMessage string
		serviceErr      error
		body            string
	}{
		{
			name:            "success response",
			expectedStatus:  http.StatusOK,
			responseMessage: "Registration success!",
			body: `{
			"username": "rawbil",
			"email": "bildadsimiyu6@gmail.com",
			"password": "@Admin123"
			}`,
		},
		{
			name:            "decode error",
			expectedStatus:  http.StatusBadRequest,
			responseMessage: "error decoding body",
			body: `{
			"userrname": "rawbil",
			"email": "bildadsimiyu6@gmail.com",
			"password": "@Admin123"
			}`,
		},
		{
			name:            "required error",
			expectedStatus:  http.StatusBadRequest,
			responseMessage: FieldsRequiredError.Error(),
			serviceErr:      FieldsRequiredError,
			body: `{
			"username": "rawbil",
			"email": "bildadsimiyu6@gmail.com",
			"password": "@Admin123"
			}`,
		},
		{
			name:            "invalid email error",
			expectedStatus:  http.StatusBadRequest,
			responseMessage: InvalidEmailError.Error(),
			serviceErr:      InvalidEmailError,
			body: `{
			"username": "rawbil",
			"email": "bildadsimiyu6@gmail.com",
			"password": "@Admin123"
			}`,
		},
		{
			name:            "invalid password error",
			expectedStatus:  http.StatusBadRequest,
			responseMessage: InvalidPasswordError.Error(),
			serviceErr:      InvalidPasswordError,
			body: `{
			"username": "rawbil",
			"email": "bildadsimiyu6@gmail.com",
			"password": "@Admin123"
			}`,
		},
		{
			name:            "conflict error",
			expectedStatus:  http.StatusConflict,
			responseMessage: UserExistsError.Error(),
			serviceErr:      UserExistsError,
			body: `{
			"username": "rawbil",
			"email": "bildadsimiyu6@gmail.com",
			"password": "@Admin123"
			}`,
		},
		{
			name:            "server error",
			expectedStatus:  http.StatusInternalServerError,
			responseMessage: "oops... server error",
			serviceErr:      errors.New("server error"),
			body: `{
			"username": "rawbil",
			"email": "bildadsimiyu6@gmail.com",
			"password": "@Admin123"
			}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := &MockAuthService{
				UserRegisterFunc: func(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error) {
					return nil, test.serviceErr
				},
			}

			h := NewHandler(mock)

			req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(test.body))

			rr := httptest.NewRecorder()

			h.UserRegister(rr, req)

			if rr.Code != test.expectedStatus {
				t.Errorf("got status %d, expected %d", rr.Code, test.expectedStatus)
			}
			if !strings.Contains(rr.Body.String(), test.responseMessage) {
				t.Errorf("expected message '%s', got '%s'", test.responseMessage, rr.Body.String())
			}
		})
	}
}

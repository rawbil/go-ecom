package auth

import (
	"context"
	"database/sql"
	"testing"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
	"github.com/rawbil/ecom2/internal/utils"
)

type MockAuthService struct {
	UserRegister func(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error)

	UserLogin            func(ctx context.Context, arg authutils.UserLoginParams) (repository.User, string, string, error)
	LogoutUser           func(ctx context.Context) error
	PasswordResetFunc    func(ctx context.Context, arg authutils.PasswordResetParams) error
	RefreshTokensFunc    func(ctx context.Context, arg authutils.RefreshTokenParam) (string, string, error)
	UpdateMyEmailFunc    func(ctx context.Context, arg authutils.UpdateEmailParams) (repository.User, error)
	UpdateMyUsernameFunc func(ctx context.Context, arg authutils.UpdateUsernameParams) (repository.User, error)
	ForgotPasswordFunc   func(ctx context.Context, arg authutils.UpdateEmailParams) error
}

// * MOCK SERVICES
func (m *MockAuthService) Register(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error) {
	return nil, nil
}

func (m *MockAuthService) Login(ctx context.Context, arg authutils.UserLoginParams) (repository.User, string, string, error) {
	return repository.User{}, "", "", nil
}

func (m *MockAuthService) Logout(ctx context.Context) error {
	return nil
}

func (m *MockAuthService) PasswordReset(ctx context.Context, arg authutils.PasswordResetParams) error {
	return nil
}

func (m *MockAuthService) RefreshTokens(ctx context.Context, arg authutils.RefreshTokenParam) (string, string, error) {
	return "", "", nil
}

func (m *MockAuthService) UpdateMyEmail(ctx context.Context, arg authutils.UpdateEmailParams) (repository.User, error) {
	return repository.User{}, nil
}

func (m *MockAuthService) UpdateMyUsername(ctx context.Context, arg authutils.UpdateUsernameParams) (repository.User, error) {
	return repository.User{}, nil
}

func (m *MockAuthService) ForgotPassword(ctx context.Context, arg authutils.UpdateEmailParams) error {
	return nil
}

// * Register Test
func TestRegister(t *testing.T) {
	utils.Slogger()

	tests := []struct{
		name string
		expectedStatus int
		serviceErr error
	}{}
}

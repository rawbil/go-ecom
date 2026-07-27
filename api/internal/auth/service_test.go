package auth

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
)

type MockRepository struct {
	ListUserFunc           func(ctx context.Context, email string) (repository.User, error)
	ListUserByIdFunc       func(ctx context.Context, userID int64) (repository.User, error)
	GetRefreshTokenFunc    func(ctx context.Context, userID int64) (repository.RefreshToken, error)
	UpdateRefreshTokenFunc func(ctx context.Context, arg repository.UpdateRefreshTokenParams) (sql.Result, error)
	CreateRefreshTokenFunc func(ctx context.Context, arg repository.CreateRefreshTokenParams) (sql.Result, error)
	DeleteRefreshTokenFunc func(ctx context.Context, userID int64) error
	UpdatePasswordFunc     func(ctx context.Context, arg repository.UpdatePasswordParams) (sql.Result, error)
	UpdateUserEmailFunc    func(ctx context.Context, arg repository.UpdateUserEmailParams) (sql.Result, error)
	CreateUserFunc         func(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error)
	UpdateUsernameFunc     func(ctx context.Context, arg repository.UpdateUsernameParams) (sql.Result, error)
	GetRoleIDFunc          func(ctx context.Context, role string) (int64, error)
	GetPermissionIDFunc    func(ctx context.Context, permission string) (int64, error)
	CreateUserRoleFunc     func(ctx context.Context, arg repository.CreateUserRoleParams) (sql.Result, error)
}

func (m *MockRepository) ListUser(ctx context.Context, email string) (repository.User, error) {
	return m.ListUserFunc(ctx, email)
}

func (m *MockRepository) ListUserById(ctx context.Context, userID int64) (repository.User, error) {
	return m.ListUserByIdFunc(ctx, userID)
}

func (m *MockRepository) GetRefreshToken(ctx context.Context, userID int64) (repository.RefreshToken, error) {
	return m.GetRefreshTokenFunc(ctx, userID)
}

func (m *MockRepository) UpdateRefreshToken(ctx context.Context, arg repository.UpdateRefreshTokenParams) (sql.Result, error) {
	return m.UpdateRefreshTokenFunc(ctx, arg)
}

func (m *MockRepository) CreateRefreshToken(ctx context.Context, arg repository.CreateRefreshTokenParams) (sql.Result, error) {
	return m.CreateRefreshTokenFunc(ctx, arg)
}

func (m *MockRepository) DeleteRefreshToken(ctx context.Context, userID int64) error {
	return m.DeleteRefreshTokenFunc(ctx, userID)
}

func (m *MockRepository) UpdatePassword(ctx context.Context, arg repository.UpdatePasswordParams) (sql.Result, error) {
	return m.UpdatePasswordFunc(ctx, arg)
}

func (m *MockRepository) UpdateUserEmail(ctx context.Context, arg repository.UpdateUserEmailParams) (sql.Result, error) {
	return m.UpdateUserEmailFunc(ctx, arg)
}

func (m *MockRepository) CreateUser(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error) {
	return m.CreateUserFunc(ctx, arg)
}

func (m *MockRepository) UpdateUsername(ctx context.Context, arg repository.UpdateUsernameParams) (sql.Result, error) {
	return m.UpdateUsernameFunc(ctx, arg)
}

func (m *MockRepository) GetRoleID(ctx context.Context, role string) (int64, error) {
	return m.GetRoleIDFunc(ctx, role)
}

func (m *MockRepository) GetPermissionID(ctx context.Context, permission string) (int64, error) {
	return m.GetPermissionIDFunc(ctx, permission)
}

func (m *MockRepository) CreateUserRole(ctx context.Context, arg repository.CreateUserRoleParams) (sql.Result, error) {
	return m.CreateUserRoleFunc(ctx, arg)
}

// ! Test Successful Registration
func TestSuccessRegisterService(t *testing.T) {
	mock := &MockRepository{
		ListUserFunc: func(ctx context.Context, email string) (repository.User, error) {
			return repository.User{}, sql.ErrNoRows
		},
		CreateUserFunc: func(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error) {
			return nil, nil
		},
	}

	svc := NewService(mock, &sql.DB{})

	_, err := svc.UserRegister(context.Background(), repository.CreateUserParams{
		Username: "rawbil",
		Email:    "bildadsimiyu@gmail.com",
		Password: "@Admin123",
	})

	if err != nil {
		t.Fatal(err)
	}

}

// ! Test User Exists Error
func TestUserExists(t *testing.T) {
	mock := &MockRepository{
		ListUserFunc: func(ctx context.Context, email string) (repository.User, error) {
			return repository.User{}, nil
		},
	}

	svc := NewService(mock, &sql.DB{})

	_, err := svc.UserRegister(context.Background(), repository.CreateUserParams{
		Username: "rawbil",
		Email:    "bildadsimiyu6@gmail.com",
		Password: "@Admin123",
	})

	if err != UserExistsError {
		t.Fatal(err)
	}

}

func TestDBFailure(t *testing.T) {
	mock := &MockRepository{
		ListUserFunc: func(ctx context.Context, email string) (repository.User, error) {
			return repository.User{}, errors.New("DB Failure")
		},
	}

	svc := NewService(mock, &sql.DB{})

	_, err := svc.UserRegister(context.Background(), repository.CreateUserParams{
		Username: "rawbil",
		Email:    "bildadsimiyu6@gmail.com",
		Password: "@Admin123",
	})

	if err == UserExistsError {
		t.Fatal(err)
	}

}

// ! Validation Tests
func TestFieldValidation(t *testing.T) {
	type Body struct {
		Username string
		Email    string
		Password string
	}
	tests := []struct {
		name        string
		returnedErr error
		body        Body
	}{
		{
			name:        "required validation",
			returnedErr: FieldsRequiredError,
			body: Body{
				Username: "rawbil",
				Email:    "bildadsimiyu6@gmail.com",
				Password: "",
			},
		},
		{
			name:        "password validation",
			returnedErr: InvalidPasswordError,
			body: Body{
				Username: "rawbil",
				Email:    "bildadsimiyu6@gmail.com",
				Password: "admin",
			},
		},
		{
			name:        "email validation",
			returnedErr: InvalidEmailError,
			body: Body{
				Username: "rawbil",
				Email:    "bildadsimiyu6@gmail",
				Password: "@Admin123",
			},
		},
		{
			name:        "username validation",
			returnedErr: UsernameLenErr,
			body: Body{
				Username: "ra",
				Email:    "bildadsimiyu6@gmail.com",
				Password: "@Admin123",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := &MockRepository{}

			svc := NewService(mock, &sql.DB{})

			_, err := svc.UserRegister(context.Background(), repository.CreateUserParams{
				Username: test.body.Username,
				Email:    test.body.Email,
				Password: test.body.Password,
			})

			if err != test.returnedErr {
				t.Fatal(err)
			}
		})
	}
}

// ! Test if password is hashed
func TestPasswordHashing(t *testing.T) {
	mock := &MockRepository{
		ListUserFunc: func(ctx context.Context, email string) (repository.User, error) {
			return repository.User{}, sql.ErrNoRows
		},
		CreateUserFunc: func(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error) {
			if arg.Password == "@Admin123" {
				t.Fatal("password wasn't hashed")
			}

			if err := authutils.ComparePasswords("@Admin123", arg.Password); err != nil {
				t.Fatal(err)
			}

			return nil, nil
		},
	}

	svc := NewService(mock, &sql.DB{})

	_, err := svc.UserRegister(context.Background(), repository.CreateUserParams{
		Username: "rawbil",
		Email:    "bildadsimiyu6@gmail.com",
		Password: "@Admin123",
	})

	if err != nil {
		t.Fatal(err)
	}
}

// ! Verify CreateUserFunc is called
func TestCreateUserCall(t *testing.T) {
	called := false
	mock := &MockRepository{
		ListUserFunc: func(ctx context.Context, email string) (repository.User, error) {
			return repository.User{}, sql.ErrNoRows
		},
		CreateUserFunc: func(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error) {
			called = true
			return nil, nil
		},
	}

	svc := NewService(mock, &sql.DB{})

	_, err := svc.UserRegister(context.Background(), repository.CreateUserParams{
		Username: "rawbil",
		Email:    "bildadsimiyu6@gmail.com",
		Password: "@Admin123",
	})

	if !called {
		t.Fatal("CreateUser not called")
	}

	if err != nil {
		t.Fatal(err)
	}
}

// * LOGIN TESTS
func TestLoginSuccess(t *testing.T) {
	// tests := []struct{
	// 	name string
	// 	validationError error
	// }{}

	mock := &MockRepository{
		ListUserFunc: func(ctx context.Context, email string) (repository.User, error) {
			hash, err := authutils.PasswordHash("@Admin123")
			if err != nil {
				return repository.User{}, err
			}
			return repository.User{Password: hash}, nil
		},
		GetRefreshTokenFunc: func(ctx context.Context, userID int64) (repository.RefreshToken, error) {
			return repository.RefreshToken{}, nil
		},
		UpdateRefreshTokenFunc: func(ctx context.Context, arg repository.UpdateRefreshTokenParams) (sql.Result, error) {
			return nil, nil
		},
	}

	svc := NewService(mock, &sql.DB{})

	_, _, _, err := svc.UserLogin(context.Background(), authutils.UserLoginParams{
		Email:    "bildadsimiyu6@gmail.com",
		Password: "@Admin123",
	})

	if err != nil {
		t.Fatal(err)
	}

}

func TestFieldErrors(t *testing.T) {
	type Body struct {
		Email    string
		Password string
	}

	tests := []struct {
		name          string
		body          Body
		expectedError error
	}{
		{
			name: "required error",
			body: Body{
				Email:    "",
				Password: "",
			},
			expectedError: FieldsRequiredError,
		},
		{
			name: "invalid email error",
			body: Body{
				Email:    "bildad",
				Password: "@Admin123",
			},
			expectedError: InvalidEmailError,
		},
		{
			name: "404 error",
			body: Body{
				Email:    "bildadsimiyu6@gmail.com",
				Password: "@Admin123",
			},
			expectedError: UserNotFoundError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := &MockRepository{
				ListUserFunc: func(ctx context.Context, email string) (repository.User, error) {
					hash, err := authutils.PasswordHash("@Admin1233")
					if err != nil {
						return repository.User{}, err
					}
					return repository.User{Password: hash}, test.expectedError
				},
			}

			svc := NewService(mock, &sql.DB{})

			_, _, _, err := svc.UserLogin(context.Background(), authutils.UserLoginParams{
				Email:    test.body.Email,
				Password: test.body.Password,
			})
			if err != test.expectedError {
				t.Fatal(err)
			}
		})
	}
}

func TestPasswordMismatchError(t *testing.T) {
	mock := &MockRepository{
		ListUserFunc: func(ctx context.Context, email string) (repository.User, error) {
			hash, err := authutils.PasswordHash("@Admin1233")
			if err != nil {
				return repository.User{}, err
			}
			return repository.User{Password: hash}, nil
		},
	}

	svc := NewService(mock, &sql.DB{})

	_, _, _, err := svc.UserLogin(context.Background(), authutils.UserLoginParams{
		Email:    "bildadsimiyu6@gmail.com",
		Password: "@Admin123",
	})
	if err != PasswordMismatchError {
		t.Fatal(err)
	}
}

func TestCreateRefreshTokenError(t *testing.T) {
	// utils.Slogger()
	// err := config.LoadEnv()
	// if err != nil {
	// 	utils.Log.Warn("No .env file found")
	// }

	mock := &MockRepository{
		ListUserFunc: func(ctx context.Context, email string) (repository.User, error) {
			hash, err := authutils.PasswordHash("@Admin123")
			if err != nil {
				return repository.User{}, err
			}
			return repository.User{
					Password: hash,
					UserID:   1,
				},
				nil
		},
		GetRefreshTokenFunc: func(ctx context.Context, userID int64) (repository.RefreshToken, error) {
			return repository.RefreshToken{}, sql.ErrNoRows
		},
		CreateRefreshTokenFunc: func(ctx context.Context, arg repository.CreateRefreshTokenParams) (sql.Result, error) {
			return nil, errors.New("error")
		},
	}

	svc := NewService(mock, &sql.DB{})

	_, _, _, err := svc.UserLogin(context.Background(), authutils.UserLoginParams{
		Email:    "bildadsimiyu6@gmail.com",
		Password: "@Admin123",
	})

	if err == nil {
		t.Fatal(err)
	}

}

func TestRefreshTokenUpdate(t *testing.T) {
	mock := &MockRepository{
		ListUserFunc: func(ctx context.Context, email string) (repository.User, error) {
			hash, err := authutils.PasswordHash("@Admin123")
			if err != nil {
				return repository.User{}, err
			}
			return repository.User{
					Password: hash,
					UserID:   1,
				},
				nil
		},
		GetRefreshTokenFunc: func(ctx context.Context, userID int64) (repository.RefreshToken, error) {
			return repository.RefreshToken{}, nil
		},
		UpdateRefreshTokenFunc: func(ctx context.Context, arg repository.UpdateRefreshTokenParams) (sql.Result, error) {
			return nil, errors.New("error")
		},
	}

	svc := NewService(mock, &sql.DB{})

	_, _, _, err := svc.UserLogin(context.Background(), authutils.UserLoginParams{
		Email:    "bildadsimiyu6@gmail.com",
		Password: "@Admin123",
	})

	if err == nil {
		t.Fatal(err)
	}

}

package auth

import (
	"context"
	"database/sql"
	"errors"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
	"github.com/rawbil/ecom2/internal/config"
)

type Service interface {
	UserRegister(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error)

	UserLogin(ctx context.Context, arg authutils.UserLoginParams) (repository.User, string, string, error)
	// UserLogout()
}

type Svc struct {
	repository repository.Queries
	db         *sql.DB
}

func NewService(repository repository.Queries, db *sql.DB) Service {
	return &Svc{
		repository: repository,
		db:         db,
	}
}

var (
	FieldsRequiredError   = errors.New("All fields are required")
	InvalidPasswordError  = errors.New("Password should have a minimum of 8 characters, have at least 1 uppercase, lowecase letter and special character")
	InvalidEmailError     = errors.New("Invalid email format")
	UserExistsError       = errors.New("User already exists")
	SqlNoRows             = errors.New("No record available")
	UserNotFoundError     = errors.New("User Not Found")
	PasswordMismatchError = errors.New("Invalid Password")
)

// ! REGISTER
func (svc *Svc) UserRegister(ctx context.Context, params repository.CreateUserParams) (sql.Result, error) {
	//& Validate fields
	if err := authutils.UserRegisterValidation(params); err != nil {
		// empty fields
		if authutils.ValidationErrorCheck("required", err) {
			return nil, FieldsRequiredError
		}
		// password error
		if authutils.ValidationErrorCheck("password_format", err) || authutils.ValidationErrorCheck("min", err) {
			return nil, InvalidPasswordError
		}
		// email error
		if authutils.ValidationErrorCheck("email", err) {
			return nil, InvalidEmailError
		}
		return nil, err
	}

	//& Ensure user does not exist before registering
	_, err := svc.repository.ListUser(ctx, params.Email)
	if err == nil {
		return nil, UserExistsError
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, UserExistsError
	}

	//& Hash Password
	hashedPassword, err := authutils.PasswordHash(params.Password)
	if err != nil {
		return nil, err
	}

	params.Password = hashedPassword

	return svc.repository.CreateUser(ctx, params)
}

// ! LOGIN
func (svc *Svc) UserLogin(ctx context.Context, arg authutils.UserLoginParams) (repository.User, string, string, error) {
	//& validate fields
	if err := authutils.UserLoginValidation(arg); err != nil {
		if authutils.ValidationErrorCheck("required", err) {
			return repository.User{}, "", "", FieldsRequiredError
		}

		if authutils.ValidationErrorCheck("email", err) {
			return repository.User{}, "", "", InvalidEmailError
		}

		return repository.User{}, "", "", err
	}

	//& Find User
	user, err := svc.repository.ListUser(ctx, arg.Email)
	if err != nil {
		return repository.User{}, "", "", UserNotFoundError
	}

	//& Compare password with stored hashed password
	if err := authutils.ComparePasswords(arg.Password, user.Password); err != nil {
		return repository.User{}, "", "", PasswordMismatchError
	}

	// & Generate authentication token
	secret := config.GetJwtConfig().JwtSecret
	if secret == "" {
		return repository.User{}, "", "", errors.New("No token secret")
	}

	token, err := authutils.GenerateAuthToken(user.UserID, []byte(secret))
	if err != nil {
		return repository.User{}, "", "", err
	}

	//& Refresh Token
	refreshToken, issued_at, expired_at, err := authutils.GenerateRefreshToken(user.UserID, []byte(secret))
	if err != nil {
		return repository.User{}, "", "", err
	}

	// & Hash Refresh Token
	hashedToken, err := authutils.PasswordHash(refreshToken)
	if err != nil {
		return repository.User{}, "", "", err
	}

	// & Save Refresh Token in Database (users and refresh_tokens transaction)

	tx, err := svc.db.Begin()
	if err != nil {
		return repository.User{}, "", "", err
	}

	defer tx.Rollback()

	qtx := svc.repository.WithTx(tx)

	rt, err := qtx.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		RefreshToken: hashedToken,
		UserID:       user.UserID,
		IssuedAt:     issued_at,
		ExpiresAt:    expired_at,
	})
	if err != nil {
		return repository.User{}, "", "", err
	}

	rt_id, err := rt.LastInsertId()
	if err != nil {
		return repository.User{}, "", "", err
	}

	if _, err := qtx.UpdateUserToken(ctx, repository.UpdateUserTokenParams{
		RefreshTokenID: sql.NullInt64{Int64: rt_id, Valid: true},
		UserID:         user.UserID,
	}); err != nil {
		return repository.User{}, "", "", err
	}

	if err := tx.Commit(); err != nil {
		return repository.User{}, "", "", err
	}

	return user, token, refreshToken, nil
}

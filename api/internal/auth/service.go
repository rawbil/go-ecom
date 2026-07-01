package auth

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"time"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
	"github.com/rawbil/ecom2/internal/config"
)

type Service interface {
	UserRegister(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error)

	UserLogin(ctx context.Context, arg authutils.UserLoginParams) (repository.User, string, string, error)
	LogoutUser(ctx context.Context) error
	PasswordReset(ctx context.Context, arg authutils.PasswordResetParams) error
	RefreshTokens(ctx context.Context, arg authutils.RefreshTokenParam) (string, string, error)
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
	AuthNotFound          = errors.New("Request User not found. Login again...")
	AuthUserNotFound      = errors.New("Authenticated User not found. Login again...")
	SimilarPasswordError  = errors.New("New password should be different from Old password")
	InvalidRefreshToken   = errors.New("Invalid refresh token. Login again...")
	TokenExpiredError     = errors.New("Refresh token expired. Login again...")
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
		if authutils.ValidationErrorCheck("password_format", err) || authutils.ValidationErrorCheck("min", err) || authutils.ValidationErrorCheck("max", err) {
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
	hashedToken := authutils.RefreshTokenHash(refreshToken)

	// & Save Refresh Token in Database
	if _, err := svc.repository.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		RefreshToken: hashedToken,
		UserID:       user.UserID,
		IssuedAt:     issued_at,
		ExpiresAt:    expired_at,
	}); err != nil {
		return repository.User{}, "", "", err
	}

	return user, token, refreshToken, nil
}

// ! LOGOUT
func (svc *Svc) LogoutUser(ctx context.Context) error {
	//& Ensure authenticated user exists from context
	user_id, ok := authutils.GetUserIDFromContext(ctx)
	if !ok {
		return AuthNotFound
	}

	//& Get User
	user, err := svc.repository.ListUserById(ctx, user_id)
	if err != nil {
		return AuthUserNotFound
	}

	//& Delete refresh_token
	if err := svc.repository.DeleteRefreshToken(ctx, user.UserID); err != nil {
		return err
	}

	return nil
}

// ! Password Reset
func (svc *Svc) PasswordReset(ctx context.Context, arg authutils.PasswordResetParams) error {
	user_id, ok := authutils.GetUserIDFromContext(ctx)
	if !ok {
		return AuthNotFound
	}

	//& Validate Fields
	if err := authutils.PasswordResetValidation(arg); err != nil {
		if authutils.ValidationErrorCheck("required", err) {
			return FieldsRequiredError
		}

		if authutils.ValidationErrorCheck("password_format", err) || authutils.ValidationErrorCheck("min", err) || authutils.ValidationErrorCheck("max", err) {
			return InvalidPasswordError
		}
		return err
	}

	//& Find user
	user, err := svc.repository.ListUserById(ctx, user_id)
	if err != nil {
		return AuthUserNotFound
	}

	//& Ensure new password is different from old password
	if err := authutils.ComparePasswords(arg.NewPassword, user.Password); err == nil {
		return SimilarPasswordError
	}

	//& Compare Passwords
	if err := authutils.ComparePasswords(arg.OldPassword, user.Password); err != nil {
		return PasswordMismatchError
	}

	//& Hash new password
	hashedPassword, err := authutils.PasswordHash(arg.NewPassword)
	if err != nil {
		return err
	}

	if _, err := svc.repository.UpdatePassword(ctx, repository.UpdatePasswordParams{
		Password: hashedPassword,
		UserID:   user.UserID,
	}); err != nil {
		return err
	}

	return nil
}

// ! Refresh Tokens
func (svc *Svc) RefreshTokens(ctx context.Context, arg authutils.RefreshTokenParam) (string, string, error) {
	user_id, ok := authutils.GetUserIDFromContext(ctx)
	if !ok {
		return "", "", AuthNotFound
	}

	//& Ensure user exists
	user, err := svc.repository.ListUserById(ctx, user_id)
	if err != nil {
		return "", "", AuthUserNotFound
	}

	//& Validate field
	if err := authutils.RefreshTokenValidation(arg); err != nil {
		if authutils.ValidationErrorCheck("required", err) {
			return "", "", FieldsRequiredError
		}
		return "", "", err
	}

	//& Find refresh token
	hashed_refresh_token, err := svc.repository.GetRefreshToken(ctx, user.UserID)
	if err != nil {
		return "", "", err
	}

	//& Validate refresh token
	claims, err := authutils.ValidateToken(arg.RefreshToken)
	if err != nil {
		return "", "", InvalidRefreshToken
	}

	//& Ensure claim id and context id are similar
	if claims.UserID != user_id {
		return "", "", InvalidRefreshToken
	}

	//& Compare refresh tokens
	client_hashed_token := authutils.RefreshTokenHash(arg.RefreshToken)

	if subtle.ConstantTimeCompare([]byte(client_hashed_token), []byte(hashed_refresh_token.RefreshToken)) != 1 {
		return "", "", InvalidRefreshToken
	}

	//& Check token expiry
	if time.Now().After(hashed_refresh_token.ExpiresAt) {
		return "", "", TokenExpiredError
	}

	//& Generate new tokens
	secret := (config.GetJwtConfig().JwtSecret)
	if secret == "" {
		return "", "", fmt.Errorf("Secret Missing")
	}

	new_refreshToken, issuedAt, expiresAt, err := authutils.GenerateRefreshToken(user.UserID, []byte(secret))
	if err != nil {
		return "", "", err
	}

	new_auth_token, err := authutils.GenerateAuthToken(user.UserID, []byte(secret))
	if err != nil {
		return "", "", err
	}

	new_hashed_token := authutils.RefreshTokenHash(new_refreshToken)

	//& Save new hashed token to DB
	if _, err := svc.repository.UpdateRefreshToken(ctx, repository.UpdateRefreshTokenParams{
		RefreshToken: new_hashed_token,
		IssuedAt:     issuedAt,
		ExpiresAt:    expiresAt,
		UserID:       user.UserID,
	}); err != nil {
		return "", "", err
	}

	return new_auth_token, new_hashed_token, nil
}

package authutils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	"github.com/rawbil/ecom2/internal/config"
)

type contextKey string

var userIDContextKey contextKey

type Claims struct {
	UserID int64 `json:"userID"`
	jwt.RegisteredClaims
}

func GenerateAuthToken(userID int64, secret []byte) (string, error) {
	expiration := time.Second * time.Duration(config.GetJwtConfig().JwtExpire)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ! AuthMiddleware
func AuthMiddleware(next http.HandlerFunc, repository repository.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Authorization Header Missing", http.StatusUnauthorized)
			return
		}

		bearer_token := strings.Split(auth, " ")
		if len(bearer_token) != 2 || bearer_token[1] != "Bearer" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		tokenString := bearer_token[1]

		//& Validate token

		claims, err := validateToken(tokenString)
		if err != nil {
			log.Printf("Failed to validate token: %v", err)
			http.Error(w, "Error validating user", http.StatusUnauthorized)
			return
		}

		//& Get user from DB with token's userID
		user, err := repository.ListUserById(r.Context(), claims.UserID)
		if err != nil {
			log.Printf("Failed to find authenticated user: %v", err)
			http.Error(w, "Error validating user", http.StatusUnauthorized)
			return
		}
		//& Get UserID from token
		ctx := context.WithValue(
			r.Context(),
			userIDContextKey,
			user.ID,
		)

		r = r.WithContext(ctx)
		next(w, r)
	}
}

// * Validate Token
func validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("Invalid Signing method")
			}
			return []byte(config.GetJwtConfig().JwtSecret), nil
		})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid Token")
	}

	return claims, nil
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDContextKey).(int64)

	return userID, ok
}

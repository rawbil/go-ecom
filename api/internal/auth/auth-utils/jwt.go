package authutils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rawbil/ecom2/internal/config"
)

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

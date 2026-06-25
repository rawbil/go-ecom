package authutils

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rawbil/ecom2/internal/config"
)

func GenerateAuthToken(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.GetJwtConfig().JwtExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

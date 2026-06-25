package authutils

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rawbil/ecom2/internal/config"
)

func GenerateAuthToken(secret []byte, user_id int) (string, error) {
	expiration := time.Second * time.Duration(config.GetJwtConfig().JwtExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"userId":    strconv.Itoa(user_id),
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	return token.SignedString(token)

}

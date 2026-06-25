package authutils

import "golang.org/x/crypto/bcrypt"

func PasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func ComparePasswords(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
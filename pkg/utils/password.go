package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	passworddigest, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(passworddigest), nil
}

func ComparePassword(passworddigest string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passworddigest), []byte(password)) == nil
}

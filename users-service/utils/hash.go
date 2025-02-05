package utils

import "golang.org/x/crypto/bcrypt"

func GenerateFromPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

func CheckPasswordHash(input, realPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(realPass), []byte(input))
	return err == nil
}

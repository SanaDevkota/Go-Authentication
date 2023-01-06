package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func Hashpassword(password string) (string, error) {
	// Hashing the password with a default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // bcrypto.DefaultCost is equal to 10
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

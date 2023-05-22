package utils

import (
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (hashedPassword string, err error) {
	bytePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytePassword), err
}

func ComparePassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func ValidEmail(email string) (isValid bool) {
	_, err := mail.ParseAddress(email)
	return err == nil
}

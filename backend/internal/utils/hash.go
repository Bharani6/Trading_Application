package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pPassword string) (string, error) {
	lBytes, lErr := bcrypt.GenerateFromPassword([]byte(pPassword), 14)
	return string(lBytes), lErr
}

func CheckPasswordHash(pPassword, pHash string) bool {
	lErr := bcrypt.CompareHashAndPassword([]byte(pHash), []byte(pPassword))
	return lErr == nil
}

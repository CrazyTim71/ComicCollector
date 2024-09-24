package crypt

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

const (
	letters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialChars = "._!?"
	numbers      = "0123456789"
)

func HashPassword(password string) (string, error) {
	// cost describes the iterations and can cause a significant delay if it is >10
	hashedPwBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	hashedPw := string(hashedPwBytes)

	return hashedPw, err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func GenerateRandomPassword(length int, includeNumber bool, includeSpecial bool) string {
	var password []byte
	var charSource string

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if includeNumber {
		charSource += numbers
	}
	if includeSpecial {
		charSource += specialChars
	}
	charSource += letters

	for i := 0; i < length; i++ {
		randNum := r.Intn(len(charSource))
		password = append(password, charSource[randNum])
	}

	return string(password)
}

package utils

import (
	"math/rand"
	"time"
)

const (
	letters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialChars = "._!?"
	numbers      = "0123456789"
)

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

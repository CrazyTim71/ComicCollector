package crypt

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// cost describes the iterations and can cause a significant delay if it is >10
	hashedPwBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	hashedPw := string(hashedPwBytes)

	return hashedPw, err
}

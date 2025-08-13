package auth

import (
	"errors"

	b "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	if len(password) > 72 {
		return "", errors.New("password is longer than the accepted limit")
	}

	hashedPW, err := b.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return "", err
	}
	return string(hashedPW), err

}

func CheckPasswordHash(p, hashp string) error {

	err := b.CompareHashAndPassword([]byte(hashp), []byte(p))

	if err != nil {
		return err
	}

	return nil
}

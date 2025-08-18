package password

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, 10)
}

func IsCorrectPassword(hashed_password, inp_password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashed_password, inp_password)
	return err == nil
}

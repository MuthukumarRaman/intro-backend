package helper

import (
	"crypto/md5"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

func CheckPasswordHashs(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// // GeneratePasswordHash - Method to generate password hash
func GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func PasswordHash(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

func ValidatePassword(password string, hashstring string) bool {

	return PasswordHash(password) == hashstring

}

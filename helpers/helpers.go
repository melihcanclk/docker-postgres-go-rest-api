package helpers

import (
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/unicode/norm"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsEmailValid(email string) (string, error) {
	mail, err := mail.ParseAddress(email)
	return mail.Address, err
}

func isIncludesNonAscii(input string) error {
	for _, r := range input {
		if r <= 127 {
			return fmt.Errorf("input contains non-ascii characters")
		}
	}
	return nil
}

func NormalizeToEnglish(input *string) error {
	normalized := norm.NFKD.String(*input)
	return isIncludesNonAscii(normalized)
}

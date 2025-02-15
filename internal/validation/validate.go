package validation

import (
	"errors"
	"unicode"
)

// Ошибки валидации
var (
	ErrInvalidUsername     = errors.New("username must be between 3 and 20 characters long")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrPasswordMinLength   = errors.New("password must be at least 8 characters long")
	ErrPasswordMaxLength   = errors.New("password must be at most 24 characters long")
	ErrPasswordNoUppercase = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLowercase = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoNumber    = errors.New("password must contain at least one number")
	ErrPasswordNoSpecial   = errors.New("password must contain at least one special character")
)

// ValidateUsername проверяет соответствие имени пользователя
func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return ErrInvalidUsername
	}
	return nil
}

// ValidatePassword проверяет соответствие пароля условиям
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordMinLength
	}
	if len(password) > 24 {
		return ErrPasswordMaxLength
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case !unicode.IsLetter(char) && !unicode.IsDigit(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUppercase
	}
	if !hasLower {
		return ErrPasswordNoLowercase
	}
	if !hasNumber {
		return ErrPasswordNoNumber
	}
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}

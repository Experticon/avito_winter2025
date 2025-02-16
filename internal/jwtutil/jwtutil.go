package jwtutil

/**
 * @package jwtutil
 * Файл jwtutil.go
 * Генерирует JWT токен на 3 часа, хранит в себе имя пользователя
 */

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(username string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   username,                                          // login пользователя
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour)), // Токен будет действовать 3 часа
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	// Создаем новый токен с использованием HMAC SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен с использованием секретного ключа
	secretKey := "secret" // Используйте свой секретный ключ
	return token.SignedString([]byte(secretKey))
}

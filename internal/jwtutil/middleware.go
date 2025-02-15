package jwtutil

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware проверяет Bearer-токен и добавляет userID в контекст запроса
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"errors": "Missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		// Проверяем, что заголовок начинается с "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"errors": "Invalid authorization header"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Разбираем и валидируем JWT-токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем метод подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			// Ключ подписи (секретный ключ)
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"errors": "Invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Извлекаем claims как map и проверяем
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			// Логируем ошибку, если не удалось привести claims
			log.Printf("Failed to parse token claims: %+v", token.Claims)
			http.Error(w, `{"errors": "Invalid token, RegisteredClaims"}`, http.StatusUnauthorized)
			return
		}

		// Извлекаем username из claims
		username, ok := claims["sub"].(string)
		if !ok || username == "" {
			http.Error(w, `{"errors": "Invalid username in token"}`, http.StatusUnauthorized)
			return
		}

		// Добавляем username в контекст запроса
		ctx := context.WithValue(r.Context(), "username", username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

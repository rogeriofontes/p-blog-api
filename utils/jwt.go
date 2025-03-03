package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Chave secreta para assinar os tokens (NÃO DEIXE ISSO NO CÓDIGO EM PRODUÇÃO)
var secretKey = []byte("super_secret_key")

// GenerateToken cria um novo token JWT para um usuário
func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expira em 24 horas
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateToken valida um token JWT e retorna os claims
func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

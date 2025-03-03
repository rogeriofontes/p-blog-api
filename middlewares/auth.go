package middlewares

import (
	"net/http"
	"strings"

	"github.com/rogeriofontes/p-blog-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware protege rotas requerendo um JWT válido
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não encontrado"})
			c.Abort()
			return
		}

		// Remover "Bearer " do token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Obter os claims do token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Erro ao processar token"})
			c.Abort()
			return
		}

		// Armazena o ID do usuário para uso posterior nos handlers
		c.Set("user_id", claims["user_id"])

		c.Next()
	}
}

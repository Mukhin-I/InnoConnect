package usecase

import (
	"fmt"
	"innoconnect/pkg/config"
	"innoconnect/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// GetUserFromToken extracts user ID and name from JWT token in the request header
func GetUserFromToken(c *gin.Context) (int64, string, error) {
	tokenStringUnparsed := c.GetHeader("Authorization")
	secret := config.GetVar("JWT_SECRET")

	const prefix = "Bearer "
	tokenString := strings.TrimPrefix(tokenStringUnparsed, prefix)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		logger.Error(err.Error())
		return 0, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", fmt.Errorf("user_id claim missing")
	}

	name, ok := claims["name"].(string)
	if !ok {
		return 0, "", fmt.Errorf("name claim missing")
	}

	return int64(userIDFloat), name, nil
}

package usecase

import (
	"fmt"
	"innoconnect/pkg/config"
	"innoconnect/pkg/logger"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func GetUserIDFromToken(tokenStringUnparsed string) (int64, error) {
	secret := config.GetVar("JWT_SECRET")
	logger.Info("JWT token " + tokenStringUnparsed)
	const prefix = "Bearer "
    tokenString := strings.TrimPrefix(tokenStringUnparsed, prefix)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// JWT numbers are decoded as float64
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id claim missing")
	}

	return int64(userIDFloat), nil
}
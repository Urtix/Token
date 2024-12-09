package generateToken

import (
	"Program/access"
	"Program/database"
	"Program/refresh"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"sync"
)

// GenerateJWT Функция для генерации JWT токена
func GenerateJWT(claims jwt.MapClaims, privateKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err) // Более информативная ошибка
	}
	return tokenString, nil
}

func GenerateAndSaveTokens(c *fiber.Ctx, userGUID uuid.UUID, db *sql.DB, privateKey []byte, tokenMap *sync.Map) (string, string, error) {
	ip := c.IP()

	// Создаем access token
	accessClaims := access.Claims(userGUID)
	accessToken, err := GenerateJWT(accessClaims, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token\": %w", err)
	}

	// Создаем refresh token
	refreshClaims := refresh.Claims(userGUID, ip)
	refreshToken, err := GenerateJWT(refreshClaims, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token\": %w", err)
	}

	// Сохраняем refresh токен в базу данных
	refreshTokenDB, err := database.CreateRefreshTokenForDB(refreshToken, userGUID, ip)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token to database\": %w", err)
	}

	if err = database.SaveRefreshTokenToDB(db, refreshTokenDB); err != nil {
		return "", "", fmt.Errorf("failed to save refresh token to database\": %w", err)
	}

	// Сохраняем связь (refresh token: access token)
	tokenMap.Store(string(refreshTokenDB.Token), accessToken)

	return accessToken, refreshToken, nil
}

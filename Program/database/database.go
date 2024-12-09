package database

import (
	"Program/token"
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// Создание таблицы c refresh токенами в бд
func CreateTableRefreshTokens(db *sql.DB) {
	_, err := db.Exec(
		"CREATE TABLE  IF NOT EXISTS tokens (" +
			"id SERIAL PRIMARY KEY," +
			"user_guid UUID NOT NULL," +
			"ip TEXT NOT NULL," +
			"token TEXT NOT NULL," +
			"expires_at TIMESTAMP WITH TIME ZONE NOT NULL," +
			"created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP" +
			");")

	if err != nil {
		panic(err)
	}
}

// Создание токена который будет храниться в базе данных
func CreateRefreshTokenForDB(refreshToken string, userGUID uuid.UUID, ip string) (*token.Token, error) {
	hasher := sha512.New()
	_, err := hasher.Write([]byte(refreshToken))
	if err != nil {
		return nil, fmt.Errorf("error hashing token: %w", err)
	}
	hashedToken := []byte(hex.EncodeToString(hasher.Sum(nil)))
	return &token.Token{
		UserGUID:  userGUID,
		IP:        ip,
		Token:     hashedToken,
		ExpiresAt: time.Now().Add(time.Hour * 24),
		CreatedAt: time.Now(),
	}, err
}

// Функция для сохранения refresh токена в базу данных
func SaveRefreshTokenToDB(db *sql.DB, token *token.Token) error {
	res, err := db.ExecContext(
		context.Background(),
		`INSERT INTO tokens (user_guid, ip, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`,
		token.UserGUID,
		token.IP,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("failed to save token: unexpected number of rows affected: %d", rowsAffected)
	}

	return nil
}

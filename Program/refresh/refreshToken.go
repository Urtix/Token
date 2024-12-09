package refresh

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"sync"
	"time"
)

func Claims(userGUID uuid.UUID, ip string) jwt.MapClaims {
	return jwt.MapClaims{
		"user_id": userGUID,
		"ip":      ip,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен истекает через 24 часа
	}
}

func DeleteRefreshToken(db *sql.DB, refreshToken []byte, tokenMap *sync.Map) error {
	_, ok := tokenMap.Load(string(refreshToken))
	if !ok {
		return fmt.Errorf("non-existent refresh token")
	}

	tokenMap.Delete(string(refreshToken))

	res, err := db.ExecContext(
		context.Background(),
		`DELETE FROM tokens WHERE token = $1;`,
		refreshToken,
	)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("failed to delete token: unexpected number of rows affected: %d", rowsAffected)
	}

	return nil
}

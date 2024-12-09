package utilits

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func IpToTokens(db *sql.DB, token []byte) (uuid.UUID, []string, error) {
	row := db.QueryRowContext(context.Background(), "SELECT user_guid FROM tokens WHERE token = $1", token)
	var userGUID uuid.UUID
	err := row.Scan(&userGUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, nil, fmt.Errorf("token not found")
		}
		return uuid.Nil, nil, fmt.Errorf("failed to scan row: %w", err)
	}

	rows, err := db.QueryContext(context.Background(), "SELECT ip FROM tokens WHERE user_guid = $1", userGUID)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("failed to query database: %w", err)
	}

	var ipList []string
	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			return uuid.Nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}
		ipList = append(ipList, ip)
	}

	if err := rows.Err(); err != nil {
		return uuid.Nil, nil, fmt.Errorf("error iterating rows: %w", err)
	}
	rows.Close()

	return userGUID, ipList, nil
}

func ParseGUID(c *fiber.Ctx) (uuid.UUID, error) {
	guid := c.Params("guid")
	userGUID, err := uuid.Parse(guid)
	if err != nil {
		return userGUID, fmt.Errorf("failed to parse GUID\": %w", err)
	}
	return userGUID, nil
}

package access

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func Claims(userGUID uuid.UUID) jwt.MapClaims {
	return jwt.MapClaims{
		"user_id": userGUID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	}
}

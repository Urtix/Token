package token

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	UserGUID  uuid.UUID `db:"user_guid"`
	IP        string    `db:"ip"`
	Token     []byte    `db:"generateToken"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

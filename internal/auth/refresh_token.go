package auth

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token	  string
	UserID 	  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
	RevokedAt time.Time
} // End RefreshToken struct

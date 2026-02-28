package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID         uuid.UUID
	IdentityID uuid.UUID
	TokenHash  string
	DeviceInfo string
	IPAddress  string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	RevokedAt  *time.Time
}

func (rt *RefreshToken) IsActive() bool {
	return rt.RevokedAt == nil && time.Now().Before(rt.ExpiresAt)
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

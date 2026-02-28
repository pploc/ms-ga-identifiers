package entity

import (
	"time"

	"github.com/google/uuid"
)

type PasswordResetToken struct {
	ID         uuid.UUID
	IdentityID uuid.UUID
	TokenHash  string
	ExpiresAt  time.Time
	UsedAt     *time.Time
	CreatedAt  time.Time
}

func (prt *PasswordResetToken) IsUsed() bool {
	return prt.UsedAt != nil
}

func (prt *PasswordResetToken) IsExpired() bool {
	return time.Now().After(prt.ExpiresAt)
}

func (prt *PasswordResetToken) IsValid() bool {
	return !prt.IsUsed() && !prt.IsExpired()
}

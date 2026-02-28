package entity

import (
	"time"

	"github.com/google/uuid"
)

type IdentityStatus string

const (
	StatusActive     IdentityStatus = "active"
	StatusLocked     IdentityStatus = "locked"
	StatusSuspended  IdentityStatus = "suspended"
	StatusUnverified IdentityStatus = "unverified"
)

type Identity struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Email         string
	PasswordHash  string
	Status        IdentityStatus
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (i *Identity) IsActive() bool {
	return i.Status == StatusActive
}

func (i *Identity) IsLocked() bool {
	return i.Status == StatusLocked
}

func (i *Identity) IsSuspended() bool {
	return i.Status == StatusSuspended
}

func (i *Identity) CanLogin() bool {
	return i.Status == StatusActive && i.EmailVerified
}

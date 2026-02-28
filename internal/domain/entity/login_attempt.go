package entity

import (
	"time"

	"github.com/google/uuid"
)

type LoginAttempt struct {
	ID          uuid.UUID
	IdentityID  *uuid.UUID
	Email       string
	IPAddress   string
	Success     bool
	AttemptedAt time.Time
}

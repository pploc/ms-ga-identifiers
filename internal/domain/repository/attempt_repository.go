package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
)

type LoginAttemptRepository interface {
	Create(ctx context.Context, attempt *entity.LoginAttempt) error
	CountRecentFailures(ctx context.Context, identityID uuid.UUID, since time.Time) (int, error)
	GetRecentByIdentityID(ctx context.Context, identityID uuid.UUID, limit int) ([]*entity.LoginAttempt, error)
}

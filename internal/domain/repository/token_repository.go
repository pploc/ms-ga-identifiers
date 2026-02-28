package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) (*entity.RefreshToken, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	GetActiveByIdentityID(ctx context.Context, identityID uuid.UUID) ([]*entity.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllByIdentityID(ctx context.Context, identityID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type PasswordResetRepository interface {
	Create(ctx context.Context, token *entity.PasswordResetToken) (*entity.PasswordResetToken, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error)
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

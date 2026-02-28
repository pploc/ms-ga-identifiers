package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
)

type IdentityRepository interface {
	Create(ctx context.Context, identity *entity.Identity) (*entity.Identity, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Identity, error)
	GetByEmail(ctx context.Context, email string) (*entity.Identity, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Identity, error)
	Update(ctx context.Context, identity *entity.Identity) (*entity.Identity, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.IdentityStatus) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	SetEmailVerified(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

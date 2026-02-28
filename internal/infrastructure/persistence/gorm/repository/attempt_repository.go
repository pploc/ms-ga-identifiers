package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
	"github.com/gym-api/ms-ga-identifier/internal/domain/repository"
	"github.com/gym-api/ms-ga-identifier/internal/infrastructure/persistence/gorm/model"
	"gorm.io/gorm"
)

type loginAttemptRepository struct {
	db *gorm.DB
}

func NewLoginAttemptRepository(db *gorm.DB) repository.LoginAttemptRepository {
	return &loginAttemptRepository{db: db}
}

func (r *loginAttemptRepository) Create(ctx context.Context, attempt *entity.LoginAttempt) error {
	m := model.EntityToLoginAttemptModel(attempt)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *loginAttemptRepository) CountRecentFailures(ctx context.Context, identityID uuid.UUID, since time.Time) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.LoginAttemptModel{}).
		Where("identity_id = ? AND success = false AND attempted_at > ?", identityID, since).
		Count(&count).Error
	return int(count), err
}

func (r *loginAttemptRepository) GetRecentByIdentityID(ctx context.Context, identityID uuid.UUID, limit int) ([]*entity.LoginAttempt, error) {
	var models []model.LoginAttemptModel
	if err := r.db.WithContext(ctx).
		Where("identity_id = ?", identityID).
		Order("attempted_at DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	attempts := make([]*entity.LoginAttempt, len(models))
	for i, m := range models {
		attempts[i] = m.ToEntity()
	}
	return attempts, nil
}

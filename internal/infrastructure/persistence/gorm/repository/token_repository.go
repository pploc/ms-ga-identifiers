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

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) repository.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) (*entity.RefreshToken, error) {
	m := model.EntityToRefreshTokenModel(token)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *refreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	var m model.RefreshTokenModel
	if err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *refreshTokenRepository) GetActiveByIdentityID(ctx context.Context, identityID uuid.UUID) ([]*entity.RefreshToken, error) {
	var models []model.RefreshTokenModel
	if err := r.db.WithContext(ctx).
		Where("identity_id = ? AND revoked_at IS NULL AND expires_at > ?", identityID, time.Now()).
		Find(&models).Error; err != nil {
		return nil, err
	}
	
	tokens := make([]*entity.RefreshToken, len(models))
	for i, m := range models {
		tokens[i] = m.ToEntity()
	}
	return tokens, nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.RefreshTokenModel{}).Where("id = ?", id).Update("revoked_at", now).Error
}

func (r *refreshTokenRepository) RevokeAllByIdentityID(ctx context.Context, identityID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.RefreshTokenModel{}).Where("identity_id = ?", identityID).Update("revoked_at", now).Error
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&model.RefreshTokenModel{}).Error
}

// PasswordResetRepository implementation
type passwordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) repository.PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(ctx context.Context, token *entity.PasswordResetToken) (*entity.PasswordResetToken, error) {
	m := model.EntityToPasswordResetModel(token)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *passwordResetRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error) {
	var m model.PasswordResetModel
	if err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *passwordResetRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.PasswordResetModel{}).Where("id = ?", id).Update("used_at", now).Error
}

func (r *passwordResetRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&model.PasswordResetModel{}).Error
}

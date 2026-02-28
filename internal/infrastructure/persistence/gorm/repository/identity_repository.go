package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
	"github.com/gym-api/ms-ga-identifier/internal/domain/repository"
	"github.com/gym-api/ms-ga-identifier/internal/infrastructure/persistence/gorm/model"
	"gorm.io/gorm"
)

type identityRepository struct {
	db *gorm.DB
}

func NewIdentityRepository(db *gorm.DB) repository.IdentityRepository {
	return &identityRepository{db: db}
}

func (r *identityRepository) Create(ctx context.Context, identity *entity.Identity) (*entity.Identity, error) {
	m := model.EntityToIdentityModel(identity)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *identityRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Identity, error) {
	var m model.IdentityModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *identityRepository) GetByEmail(ctx context.Context, email string) (*entity.Identity, error) {
	var m model.IdentityModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *identityRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Identity, error) {
	var m model.IdentityModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *identityRepository) Update(ctx context.Context, identity *entity.Identity) (*entity.Identity, error) {
	m := model.EntityToIdentityModel(identity)
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *identityRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.IdentityStatus) error {
	return r.db.WithContext(ctx).Model(&model.IdentityModel{}).Where("id = ?", id).Update("status", status).Error
}

func (r *identityRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	return r.db.WithContext(ctx).Model(&model.IdentityModel{}).Where("id = ?", id).Update("password_hash", passwordHash).Error
}

func (r *identityRepository) SetEmailVerified(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.IdentityModel{}).Where("id = ?", id).Update("email_verified", true).Error
}

func (r *identityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.IdentityModel{}, "id = ?", id).Error
}

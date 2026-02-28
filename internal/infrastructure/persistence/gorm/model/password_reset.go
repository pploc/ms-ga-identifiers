package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
)

type PasswordResetModel struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	IdentityID uuid.UUID  `gorm:"type:uuid;not null;index"`
	TokenHash  string     `gorm:"type:varchar(64);uniqueIndex;not null"`
	ExpiresAt  time.Time  `gorm:"not null;index"`
	UsedAt     *time.Time `gorm:"index"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
}

func (PasswordResetModel) TableName() string {
	return "password_reset_tokens"
}

func (m *PasswordResetModel) ToEntity() *entity.PasswordResetToken {
	return &entity.PasswordResetToken{
		ID:         m.ID,
		IdentityID: m.IdentityID,
		TokenHash:  m.TokenHash,
		ExpiresAt:  m.ExpiresAt,
		UsedAt:     m.UsedAt,
		CreatedAt:  m.CreatedAt,
	}
}

func EntityToPasswordResetModel(e *entity.PasswordResetToken) *PasswordResetModel {
	return &PasswordResetModel{
		ID:         e.ID,
		IdentityID: e.IdentityID,
		TokenHash:  e.TokenHash,
		ExpiresAt:  e.ExpiresAt,
		UsedAt:     e.UsedAt,
		CreatedAt:  e.CreatedAt,
	}
}

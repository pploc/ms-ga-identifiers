package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
)

type RefreshTokenModel struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	IdentityID uuid.UUID `gorm:"type:uuid;not null;index:idx_refresh_tokens_identity_id"`
	TokenHash  string     `gorm:"type:varchar(64);uniqueIndex;not null"`
	DeviceInfo string     `gorm:"type:varchar(255)"`
	IPAddress  string     `gorm:"type:varchar(45)"`
	ExpiresAt  time.Time  `gorm:"not null;index:idx_refresh_tokens_expires_at"`
	RevokedAt  *time.Time `gorm:"index"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
}

func (RefreshTokenModel) TableName() string {
	return "refresh_tokens"
}

func (m *RefreshTokenModel) ToEntity() *entity.RefreshToken {
	return &entity.RefreshToken{
		ID:         m.ID,
		IdentityID: m.IdentityID,
		TokenHash:  m.TokenHash,
		DeviceInfo: m.DeviceInfo,
		IPAddress:  m.IPAddress,
		ExpiresAt:  m.ExpiresAt,
		CreatedAt:  m.CreatedAt,
		RevokedAt:  m.RevokedAt,
	}
}

func EntityToRefreshTokenModel(e *entity.RefreshToken) *RefreshTokenModel {
	return &RefreshTokenModel{
		ID:         e.ID,
		IdentityID: e.IdentityID,
		TokenHash:  e.TokenHash,
		DeviceInfo: e.DeviceInfo,
		IPAddress:  e.IPAddress,
		ExpiresAt:  e.ExpiresAt,
		CreatedAt:  e.CreatedAt,
		RevokedAt:  e.RevokedAt,
	}
}

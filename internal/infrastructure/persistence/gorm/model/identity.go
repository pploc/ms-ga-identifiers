package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
)

type IdentityModel struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null"`
	Email         string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash  string         `gorm:"type:varchar(255);not null"`
	Status        string         `gorm:"type:varchar(20);default:unverified"`
	EmailVerified bool           `gorm:"default:false"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
}

func (IdentityModel) TableName() string {
	return "identities"
}

func (m *IdentityModel) ToEntity() *entity.Identity {
	return &entity.Identity{
		ID:            m.ID,
		UserID:        m.UserID,
		Email:         m.Email,
		PasswordHash:  m.PasswordHash,
		Status:        entity.IdentityStatus(m.Status),
		EmailVerified: m.EmailVerified,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func EntityToIdentityModel(e *entity.Identity) *IdentityModel {
	return &IdentityModel{
		ID:            e.ID,
		UserID:        e.UserID,
		Email:         e.Email,
		PasswordHash:  e.PasswordHash,
		Status:        string(e.Status),
		EmailVerified: e.EmailVerified,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

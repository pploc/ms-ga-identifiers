package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
)

type LoginAttemptModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	IdentityID  *uuid.UUID `gorm:"type:uuid;index:idx_login_attempts_identity_id"`
	Email       string     `gorm:"type:varchar(255);not null;index:idx_login_attempts_email"`
	IPAddress   string     `gorm:"type:varchar(45)"`
	Success     bool       `gorm:"not null;default:false"`
	AttemptedAt time.Time  `gorm:"not null;index:idx_login_attempts_attempted_at;autoCreateTime"`
}

func (LoginAttemptModel) TableName() string {
	return "login_attempts"
}

func (m *LoginAttemptModel) ToEntity() *entity.LoginAttempt {
	return &entity.LoginAttempt{
		ID:          m.ID,
		IdentityID:  m.IdentityID,
		Email:       m.Email,
		IPAddress:   m.IPAddress,
		Success:     m.Success,
		AttemptedAt: m.AttemptedAt,
	}
}

func EntityToLoginAttemptModel(e *entity.LoginAttempt) *LoginAttemptModel {
	return &LoginAttemptModel{
		ID:          e.ID,
		IdentityID:  e.IdentityID,
		Email:       e.Email,
		IPAddress:   e.IPAddress,
		Success:     e.Success,
		AttemptedAt: e.AttemptedAt,
	}
}

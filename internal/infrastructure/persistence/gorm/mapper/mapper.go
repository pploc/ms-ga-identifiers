package mapper

import (
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
	"github.com/gym-api/ms-ga-identifier/internal/infrastructure/persistence/gorm/model"
)

// IdentityModelToEntity converts GORM model to domain entity
func IdentityModelToEntity(m *model.IdentityModel) *entity.Identity {
	return m.ToEntity()
}

// EntityToIdentityModel converts domain entity to GORM model
func EntityToIdentityModel(e *entity.Identity) *model.IdentityModel {
	return model.EntityToIdentityModel(e)
}

// RefreshTokenModelToEntity converts GORM model to domain entity
func RefreshTokenModelToEntity(m *model.RefreshTokenModel) *entity.RefreshToken {
	return m.ToEntity()
}

// EntityToRefreshTokenModel converts domain entity to GORM model
func EntityToRefreshTokenModel(e *entity.RefreshToken) *model.RefreshTokenModel {
	return model.EntityToRefreshTokenModel(e)
}

// LoginAttemptModelToEntity converts GORM model to domain entity
func LoginAttemptModelToEntity(m *model.LoginAttemptModel) *entity.LoginAttempt {
	return m.ToEntity()
}

// EntityToLoginAttemptModel converts domain entity to GORM model
func EntityToLoginAttemptModel(e *entity.LoginAttempt) *model.LoginAttemptModel {
	return model.EntityToLoginAttemptModel(e)
}

// PasswordResetModelToEntity converts GORM model to domain entity
func PasswordResetModelToEntity(m *model.PasswordResetModel) *entity.PasswordResetToken {
	return m.ToEntity()
}

// EntityToPasswordResetModel converts domain entity to GORM model
func EntityToPasswordResetModel(e *entity.PasswordResetToken) *model.PasswordResetModel {
	return model.EntityToPasswordResetModel(e)
}

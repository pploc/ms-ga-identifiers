package converter

import (
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
)

// IdentityToEntity converts API model to domain entity
func IdentityToEntity(id *entity.Identity) *entity.Identity {
	return id
}

// EntityToIdentity converts domain entity to API model
func EntityToIdentity(e *entity.Identity) *entity.Identity {
	return e
}

// RefreshTokenToEntity converts API model to domain entity
func RefreshTokenToEntity(rt *entity.RefreshToken) *entity.RefreshToken {
	return rt
}

// EntityToRefreshToken converts domain entity to API model
func EntityToRefreshToken(e *entity.RefreshToken) *entity.RefreshToken {
	return e
}

// LoginAttemptToEntity converts API model to domain entity
func LoginAttemptToEntity(la *entity.LoginAttempt) *entity.LoginAttempt {
	return la
}

// EntityToLoginAttempt converts domain entity to API model
func EntityToLoginAttempt(e *entity.LoginAttempt) *entity.LoginAttempt {
	return e
}

// PasswordResetToEntity converts API model to domain entity
func PasswordResetToEntity(pr *entity.PasswordResetToken) *entity.PasswordResetToken {
	return pr
}

// EntityToPasswordReset converts domain entity to API model
func EntityToPasswordReset(e *entity.PasswordResetToken) *entity.PasswordResetToken {
	return e
}

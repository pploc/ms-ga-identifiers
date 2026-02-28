package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
	"github.com/gym-api/ms-ga-identifier/internal/domain/repository"
	"github.com/gym-api/ms-ga-identifier/pkg/utils"
	"gorm.io/gorm"
)

type PasswordService struct {
	identityRepo repository.IdentityRepository
	passwordRepo repository.PasswordResetRepository
}

func NewPasswordService(
	identityRepo repository.IdentityRepository,
	passwordRepo repository.PasswordResetRepository,
) *PasswordService {
	return &PasswordService{
		identityRepo: identityRepo,
		passwordRepo: passwordRepo,
	}
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

func (s *PasswordService) ForgotPassword(ctx context.Context, email string) (*ForgotPasswordResponse, error) {
	// Always return success to prevent email enumeration
	// But only proceed if email exists
	identity, err := s.identityRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Email doesn't exist, but we still return success
			return &ForgotPasswordResponse{
				Message: "If the email exists, a reset link has been sent.",
			}, nil
		}
		return nil, err
	}

	// Generate reset token
	token := generateToken()
	tokenHash := utils.HashToken(token)

	resetToken := &entity.PasswordResetToken{
		ID:         uuid.New(),
		IdentityID: identity.ID,
		TokenHash:  tokenHash,
		ExpiresAt:  time.Now().Add(1 * time.Hour), // Token valid for 1 hour
		CreatedAt:  time.Now(),
	}

	_, err = s.passwordRepo.Create(ctx, resetToken)
	if err != nil {
		return nil, err
	}

	// In production, send email with reset link
	// For now, just log the token (in development)
	utils.Info("Password reset token generated", utils.String("token", token))

	return &ForgotPasswordResponse{
		Message: "If the email exists, a reset link has been sent.",
	}, nil
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

func (s *PasswordService) ResetPassword(ctx context.Context, token, newPassword string) (*ResetPasswordResponse, error) {
	// Hash the provided token
	tokenHash := utils.HashToken(token)

	// Look up in DB
	resetToken, err := s.passwordRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired reset token")
		}
		return nil, err
	}

	// Check token validity
	if !resetToken.IsValid() {
		return nil, errors.New("reset token is expired or already used")
	}

	// Hash new password
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}

	// Update password
	err = s.identityRepo.UpdatePassword(ctx, resetToken.IdentityID, passwordHash)
	if err != nil {
		return nil, err
	}

	// Mark token as used
	err = s.passwordRepo.MarkAsUsed(ctx, resetToken.ID)
	if err != nil {
		return nil, err
	}

	// Revoke all existing refresh tokens for security
	// This would require access to the token repository

	return &ResetPasswordResponse{
		Message: "Password reset successful.",
	}, nil
}

package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
	"github.com/gym-api/ms-ga-identifier/internal/domain/repository"
	"github.com/gym-api/ms-ga-identifier/internal/infrastructure/external"
	"github.com/gym-api/ms-ga-identifier/pkg/config"
	"github.com/gym-api/ms-ga-identifier/pkg/utils"
	"gorm.io/gorm"
)

type TokenService struct {
	identityRepo repository.IdentityRepository
	tokenRepo    repository.RefreshTokenRepository
	authClient   *external.AuthClient
	jwtUtil     *utils.JWTUtil
	cfg         *config.Config
}

func NewTokenService(
	identityRepo repository.IdentityRepository,
	tokenRepo repository.RefreshTokenRepository,
	authClient *external.AuthClient,
	jwtUtil *utils.JWTUtil,
	cfg *config.Config,
) *TokenService {
	return &TokenService{
		identityRepo: identityRepo,
		tokenRepo:    tokenRepo,
		authClient:   authClient,
		jwtUtil:     jwtUtil,
		cfg:         cfg,
	}
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (s *TokenService) RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	// Hash the provided refresh token
	refreshTokenHash := utils.HashToken(refreshToken)

	// Look up in DB
	tokenEntity, err := s.tokenRepo.GetByTokenHash(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired refresh token")
		}
		return nil, err
	}

	// Check not expired and not revoked
	if !tokenEntity.IsActive() {
		return nil, errors.New("refresh token is expired or revoked")
	}

	// Get identity
	identity, err := s.identityRepo.GetByID(ctx, tokenEntity.IdentityID)
	if err != nil {
		return nil, err
	}

	// Get fresh roles and permissions from auth service
	roles, permissions, err := s.authClient.ExtractRolesAndPermissions(identity.UserID)
	if err != nil {
		utils.Errorf("Failed to get roles and permissions during token refresh", utils.ErrorField(err.Error()))
	}

	// Issue new access token
	accessToken, err := s.jwtUtil.GenerateToken(identity.UserID.String(), identity.Email, roles, permissions)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(s.cfg.JWT.ExpirationTime.Seconds()),
	}, nil
}

func (s *TokenService) GenerateRefreshToken(ctx context.Context, identityID uuid.UUID, deviceInfo, ipAddress string) (string, error) {
	// Generate refresh token
	refreshToken := generateToken()
	refreshTokenHash := utils.HashToken(refreshToken)

	refreshTokenEntity := &entity.RefreshToken{
		ID:         uuid.New(),
		IdentityID: identityID,
		TokenHash:  refreshTokenHash,
		DeviceInfo: deviceInfo,
		IPAddress:  ipAddress,
		ExpiresAt:  time.Now().Add(s.cfg.JWT.RefreshDuration),
		CreatedAt:  time.Now(),
	}

	_, err := s.tokenRepo.Create(ctx, refreshTokenEntity)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

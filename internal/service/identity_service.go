package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gym-api/ms-ga-identifier/internal/domain/entity"
	"github.com/gym-api/ms-ga-identifier/internal/domain/repository"
	"github.com/gym-api/ms-ga-identifier/internal/infrastructure/external"
	"github.com/gym-api/ms-ga-identifier/internal/infrastructure/messaging"
	"github.com/gym-api/ms-ga-identifier/pkg/config"
	"github.com/gym-api/ms-ga-identifier/pkg/utils"
	"gorm.io/gorm"
)

type IdentityService struct {
	identityRepo  repository.IdentityRepository
	tokenRepo     repository.RefreshTokenRepository
	attemptRepo   repository.LoginAttemptRepository
	passwordRepo  repository.PasswordResetRepository
	authClient    *external.AuthClient
	kafkaProducer *messaging.KafkaProducer
	jwtUtil       *utils.JWTUtil
	cfg           *config.Config
}

func NewIdentityService(
	identityRepo repository.IdentityRepository,
	tokenRepo repository.RefreshTokenRepository,
	attemptRepo repository.LoginAttemptRepository,
	passwordRepo repository.PasswordResetRepository,
	authClient *external.AuthClient,
	kafkaProducer *messaging.KafkaProducer,
	jwtUtil *utils.JWTUtil,
	cfg *config.Config,
) *IdentityService {
	return &IdentityService{
		identityRepo:  identityRepo,
		tokenRepo:     tokenRepo,
		attemptRepo:   attemptRepo,
		passwordRepo:  passwordRepo,
		authClient:    authClient,
		kafkaProducer: kafkaProducer,
		jwtUtil:       jwtUtil,
		cfg:           cfg,
	}
}

type RegisterRequest struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type RegisterResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Email   string   `json:"email"`
	Message string   `json:"message"`
}

func (s *IdentityService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// Check if email already exists
	existing, err := s.identityRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Generate user ID
	userID := uuid.New()

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create identity
	identity := &entity.Identity{
		ID:            uuid.New(),
		UserID:        userID,
		Email:         req.Email,
		PasswordHash:  passwordHash,
		Status:        entity.StatusUnverified,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = s.identityRepo.Create(ctx, identity)
	if err != nil {
		return nil, err
	}

	// Publish event
	if s.kafkaProducer != nil {
		s.kafkaProducer.PublishIdentityRegistered(ctx, userID.String(), req.Email)
	}

	return &RegisterResponse{
		UserID:  userID,
		Email:   req.Email,
		Message: "Registration successful. Please verify your email.",
	}, nil
}

type LoginRequest struct {
	Email      string
	Password   string
	DeviceInfo string
	IPAddress  string
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func (s *IdentityService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Find identity by email
	identity, err := s.identityRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Check identity status
	if identity.Status == entity.StatusLocked || identity.Status == entity.StatusSuspended {
		return nil, errors.New("account locked or suspended")
	}

	// Check password
	if !utils.CheckPassword(req.Password, identity.PasswordHash) {
		// Record failed attempt
		s.recordLoginAttempt(ctx, identity, req.Email, req.IPAddress, false)
		return nil, errors.New("invalid credentials")
	}

	// Record successful attempt
	s.recordLoginAttempt(ctx, identity, req.Email, req.IPAddress, true)

	// Get roles and permissions from auth service
	roles, permissions, err := s.authClient.ExtractRolesAndPermissions(identity.UserID)
	if err != nil {
		utils.Errorf("Failed to get roles and permissions", utils.ErrorField(err.Error()))
	}

	// Generate JWT
	accessToken, err := s.jwtUtil.GenerateToken(identity.UserID.String(), identity.Email, roles, permissions)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken := generateToken()
	refreshTokenHash := utils.HashToken(refreshToken)

	refreshTokenEntity := &entity.RefreshToken{
		ID:         uuid.New(),
		IdentityID: identity.ID,
		TokenHash:  refreshTokenHash,
		DeviceInfo: req.DeviceInfo,
		IPAddress:  req.IPAddress,
		ExpiresAt:  time.Now().Add(s.cfg.JWT.RefreshDuration),
		CreatedAt:  time.Now(),
	}

	_, err = s.tokenRepo.Create(ctx, refreshTokenEntity)
	if err != nil {
		return nil, err
	}

	// Publish login event
	if s.kafkaProducer != nil {
		s.kafkaProducer.PublishIdentityLoggedIn(ctx, identity.UserID.String(), identity.Email, map[string]interface{}{
			"device_info": req.DeviceInfo,
			"ip_address":  req.IPAddress,
		})
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.cfg.JWT.ExpirationTime.Seconds()),
	}, nil
}

func (s *IdentityService) Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	// Revoke all refresh tokens
	err := s.tokenRepo.RevokeAllByIdentityID(ctx, userID)
	if err != nil {
		return err
	}

	// Get identity for publishing event
	identity, err := s.identityRepo.GetByUserID(ctx, userID)
	if err == nil && s.kafkaProducer != nil {
		s.kafkaProducer.PublishIdentityLoggedOut(ctx, userID.String(), identity.Email)
	}

	return nil
}

func (s *IdentityService) recordLoginAttempt(ctx context.Context, identity *entity.Identity, email, ipAddress string, success bool) {
	attempt := &entity.LoginAttempt{
		ID:          uuid.New(),
		IdentityID:  &identity.ID,
		Email:       email,
		IPAddress:   ipAddress,
		Success:     success,
		AttemptedAt: time.Now(),
	}
	s.attemptRepo.Create(ctx, attempt)
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

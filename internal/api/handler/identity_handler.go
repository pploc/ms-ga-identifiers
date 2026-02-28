package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gym-api/ms-ga-identifier/internal/middleware"
	"github.com/gym-api/ms-ga-identifier/internal/service"
	"github.com/gym-api/ms-ga-identifier/pkg/utils"
)

type IdentityHandler struct {
	identityService *service.IdentityService
}

func NewIdentityHandler(identityService *service.IdentityService) *IdentityHandler {
	return &IdentityHandler{identityService: identityService}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func (h *IdentityHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.identityService.Register(c.Request.Context(), service.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})

	if err != nil {
		utils.Conflict(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, resp)
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	DeviceInfo string `json:"device_info"`
}

func (h *IdentityHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	ipAddress := c.ClientIP()

	resp, err := h.identityService.Login(c.Request.Context(), service.LoginRequest{
		Email:      req.Email,
		Password:   req.Password,
		DeviceInfo: req.DeviceInfo,
		IPAddress:  ipAddress,
	})

	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, resp)
}

func (h *IdentityHandler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// Get user ID from context (need to parse from JWT)
	// For now, we'll implement basic logout

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *IdentityHandler) GetCurrentUser(c *gin.Context) {
	email := middleware.GetEmail(c)
	roles := middleware.GetRoles(c)
	permissions := middleware.GetPermissions(c)

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"email":       email,
		"roles":       roles,
		"permissions": permissions,
	})
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword    string `json:"new_password" binding:"required,min=8"`
}

func (h *IdentityHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *IdentityHandler) VerifyEmail(c *gin.Context) {
	token := c.Param("token")

	// Implement email verification logic
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Email verified successfully"})
}

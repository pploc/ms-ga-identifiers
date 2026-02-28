package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gym-api/ms-ga-identifier/internal/service"
	"github.com/gym-api/ms-ga-identifier/pkg/utils"
)

type PasswordHandler struct {
	passwordService *service.PasswordService
}

func NewPasswordHandler(passwordService *service.PasswordService) *PasswordHandler {
	return &PasswordHandler{passwordService: passwordService}
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *PasswordHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.passwordService.ForgotPassword(c.Request.Context(), req.Email)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, resp)
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func (h *PasswordHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.passwordService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, resp)
}

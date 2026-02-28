package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gym-api/ms-ga-identifier/internal/api/handler"
	"github.com/gym-api/ms-ga-identifier/internal/middleware"
	"github.com/gym-api/ms-ga-identifier/pkg/config"
	"github.com/gym-api/ms-ga-identifier/pkg/utils"
)

type Router struct {
	engine           *gin.Engine
	identityHandler  *handler.IdentityHandler
	tokenHandler     *handler.TokenHandler
	passwordHandler  *handler.PasswordHandler
	authMiddleware   *middleware.AuthMiddleware
	cfg              *config.Config
}

func NewRouter(
	identityHandler *handler.IdentityHandler,
	tokenHandler *handler.TokenHandler,
	passwordHandler *handler.PasswordHandler,
	authMiddleware *middleware.AuthMiddleware,
	cfg *config.Config,
) *gin.Engine {
	r := &Router{
		engine:          gin.Default(),
		identityHandler: identityHandler,
		tokenHandler:    tokenHandler,
		passwordHandler: passwordHandler,
		authMiddleware:  authMiddleware,
		cfg:            cfg,
	}

	r.setupRoutes()
	return r.engine
}

func (r *Router) setupRoutes() {
	// Health check endpoints
	r.engine.GET("/health", func(c *gin.Context) {
		utils.SuccessResponse(c, http.StatusOK, gin.H{"status": "healthy"})
	})

	r.engine.GET("/ready", func(c *gin.Context) {
		utils.SuccessResponse(c, http.StatusOK, gin.H{"status": "ready"})
	})

	// API routes
	api := r.engine.Group("/api/v1")
	{
		// Public endpoints
		public := api.Group("")
		{
			public.POST("/register", r.identityHandler.Register)
			public.POST("/login", r.identityHandler.Login)
			public.POST("/refresh", r.tokenHandler.RefreshToken)
			public.POST("/forgot-password", r.passwordHandler.ForgotPassword)
			public.POST("/reset-password", r.passwordHandler.ResetPassword)
			public.GET("/verify-email/:token", r.identityHandler.VerifyEmail)
		}

		// Protected endpoints
		protected := api.Group("")
		protected.Use(r.authMiddleware.RequireAuth())
		{
			protected.POST("/logout", r.identityHandler.Logout)
			protected.GET("/me", r.identityHandler.GetCurrentUser)
			protected.POST("/change-password", r.identityHandler.ChangePassword)
		}
	}
}

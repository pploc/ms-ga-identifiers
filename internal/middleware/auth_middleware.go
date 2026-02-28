package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gym-api/ms-ga-identifier/pkg/utils"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "user_id"
	EmailKey            = "email"
	RolesKey            = "roles"
	PermissionsKey      = "permissions"
)

type AuthMiddleware struct {
	jwtUtil *utils.JWTUtil
}

func NewAuthMiddleware(jwtUtil *utils.JWTUtil) *AuthMiddleware {
	return &AuthMiddleware{jwtUtil: jwtUtil}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			utils.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			utils.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)

		claims, err := m.jwtUtil.ValidateToken(tokenString)
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(EmailKey, claims.Email)
		c.Set(RolesKey, claims.Roles)
		c.Set(PermissionsKey, claims.Permissions)

		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get(UserIDKey); exists {
		return userID.(string)
	}
	return ""
}

func GetEmail(c *gin.Context) string {
	if email, exists := c.Get(EmailKey); exists {
		return email.(string)
	}
	return ""
}

func GetRoles(c *gin.Context) []string {
	if roles, exists := c.Get(RolesKey); exists {
		return roles.([]string)
	}
	return []string{}
}

func GetPermissions(c *gin.Context) []string {
	if permissions, exists := c.Get(PermissionsKey); exists {
		return permissions.([]string)
	}
	return []string{}
}

func HasPermission(c *gin.Context, permission string) bool {
	permissions := GetPermissions(c)
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func HasRole(c *gin.Context, role string) bool {
	roles := GetRoles(c)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/service"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/infrastructure/external"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService  *service.AuthDomainService
	cacheService *external.CacheService
}

func NewAuthMiddleware(authService *service.AuthDomainService, cacheService *external.CacheService) *AuthMiddleware {
	return &AuthMiddleware{
		authService:  authService,
		cacheService: cacheService,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		if m.cacheService != nil {
			blacklisted, err := m.cacheService.IsTokenBlacklisted(c.Request.Context(), token)
			if err == nil && blacklisted {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
				c.Abort()
				return
			}
		}

		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		c.Set("token", token)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("user_roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User roles not found"})
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user roles"})
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range userRoles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) RequireAnyRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, err := m.getUserRoles(c)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !m.hasAnyRole(userRoles, requiredRoles) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) getUserRoles(c *gin.Context) ([]string, error) {
	roles, exists := c.Get("user_roles")
	if !exists {
		return nil, fmt.Errorf("user roles not found")
	}

	userRoles, ok := roles.([]string)
	if !ok {
		return nil, fmt.Errorf("invalid user roles")
	}

	return userRoles, nil
}

func (m *AuthMiddleware) hasAnyRole(userRoles []string, requiredRoles []string) bool {
	for _, userRole := range userRoles {
		for _, requiredRole := range requiredRoles {
			if userRole == requiredRole {
				return true
			}
		}
	}
	return false
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			c.Next()
			return
		}

		if m.cacheService != nil {
			blacklisted, err := m.cacheService.IsTokenBlacklisted(c.Request.Context(), token)
			if err == nil && blacklisted {
				c.Next()
				return
			}
		}

		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		c.Set("token", token)

		c.Next()
	}
}

func (m *AuthMiddleware) extractToken(c *gin.Context) string {
	if token := m.extractFromHeader(c); token != "" {
		return token
	}

	if token := c.Query("token"); token != "" {
		return token
	}

	if token := m.extractFromCookie(c); token != "" {
		return token
	}

	return ""
}

func (m *AuthMiddleware) extractFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func (m *AuthMiddleware) extractFromCookie(c *gin.Context) string {
	cookie, err := c.Cookie("access_token")
	if err != nil {
		return ""
	}
	return cookie
}

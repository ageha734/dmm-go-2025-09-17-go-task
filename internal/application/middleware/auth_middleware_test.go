package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequireAuthNoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMiddleware := middleware.NewAuthMiddleware(nil, nil)

	router := gin.New()
	router.Use(authMiddleware.RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization token required")
}

func TestRequireRoleValidRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMiddleware := middleware.NewAuthMiddleware(nil, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_roles", []string{"admin", "user"})
		c.Next()
	})
	router.Use(authMiddleware.RequireRole("admin"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRoleInvalidRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMiddleware := middleware.NewAuthMiddleware(nil, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_roles", []string{"user"})
		c.Next()
	})
	router.Use(authMiddleware.RequireRole("admin"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Insufficient permissions")
}

func TestRequireRoleNoRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMiddleware := middleware.NewAuthMiddleware(nil, nil)

	router := gin.New()
	router.Use(authMiddleware.RequireRole("admin"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "User roles not found")
}

func TestRequireAnyRoleValidRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMiddleware := middleware.NewAuthMiddleware(nil, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_roles", []string{"moderator"})
		c.Next()
	})
	router.Use(authMiddleware.RequireAnyRole("admin", "moderator"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireAnyRoleInvalidRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMiddleware := middleware.NewAuthMiddleware(nil, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_roles", []string{"user"})
		c.Next()
	})
	router.Use(authMiddleware.RequireAnyRole("admin", "moderator"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Insufficient permissions")
}

func TestOptionalAuthNoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMiddleware := middleware.NewAuthMiddleware(nil, nil)

	router := gin.New()
	router.Use(authMiddleware.OptionalAuth())
	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if exists {
			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "anonymous"})
		}
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "anonymous")
}

func TestInvalidAuthorizationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMiddleware := middleware.NewAuthMiddleware(nil, nil)

	router := gin.New()
	router.Use(authMiddleware.RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	tests := []struct {
		name   string
		header string
	}{
		{"No Bearer prefix", "invalid-token"},
		{"Wrong prefix", "Basic invalid-token"},
		{"Empty token", "Bearer "},
		{"Only Bearer", "Bearer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.header)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestRequireRoleInvalidRoleType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMiddleware := middleware.NewAuthMiddleware(nil, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_roles", "invalid-type")
		c.Next()
	})
	router.Use(authMiddleware.RequireRole("admin"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid user roles")
}

func TestRequireAnyRoleNoRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMiddleware := middleware.NewAuthMiddleware(nil, nil)

	router := gin.New()
	router.Use(authMiddleware.RequireAnyRole("admin", "moderator"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "user roles not found")
}

func TestRequireAnyRoleInvalidRoleType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMiddleware := middleware.NewAuthMiddleware(nil, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_roles", 123)
		c.Next()
	})
	router.Use(authMiddleware.RequireAnyRole("admin", "moderator"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "invalid user roles")
}

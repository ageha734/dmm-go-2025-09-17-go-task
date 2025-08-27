package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitByIPNilCacheService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimitMiddleware := middleware.NewRateLimitMiddleware(nil)

	router := gin.New()
	router.Use(rateLimitMiddleware.RateLimitByIP(10, time.Minute))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimitByUserNoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimitMiddleware := middleware.NewRateLimitMiddleware(nil)

	router := gin.New()
	router.Use(rateLimitMiddleware.RateLimitByUser(100, time.Hour))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimitByUserWithUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimitMiddleware := middleware.NewRateLimitMiddleware(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", uint(123))
		c.Next()
	})
	router.Use(rateLimitMiddleware.RateLimitByUser(100, time.Hour))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimitByEndpointWithinLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimitMiddleware := middleware.NewRateLimitMiddleware(nil)

	router := gin.New()
	router.Use(rateLimitMiddleware.RateLimitByEndpoint("/api/login", 5, time.Minute*15))
	router.POST("/api/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "login success"})
	})

	req, _ := http.NewRequest("POST", "/api/login", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCheckBlacklistNilCacheService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimitMiddleware := middleware.NewRateLimitMiddleware(nil)

	router := gin.New()
	router.Use(rateLimitMiddleware.CheckBlacklist())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimitHeadersFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimitMiddleware := middleware.NewRateLimitMiddleware(nil)

	router := gin.New()
	router.Use(rateLimitMiddleware.RateLimitByIP(10, time.Minute))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	limitHeader := w.Header().Get("X-RateLimit-Limit")
	remainingHeader := w.Header().Get("X-RateLimit-Remaining")
	resetHeader := w.Header().Get("X-RateLimit-Reset")

	assert.NotEmpty(t, limitHeader)
	assert.NotEmpty(t, remainingHeader)
	assert.NotEmpty(t, resetHeader)
}

func TestMultipleMiddlewareChaining(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimitMiddleware := middleware.NewRateLimitMiddleware(nil)

	router := gin.New()
	router.Use(rateLimitMiddleware.CheckBlacklist())
	router.Use(rateLimitMiddleware.RateLimitByIP(10, time.Minute))
	router.Use(rateLimitMiddleware.RateLimitByEndpoint("/test", 5, time.Minute))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimitByUserDifferentUserTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name   string
		userID interface{}
	}{
		{"uint user ID", uint(123)},
		{"int user ID", 456},
		{"string user ID", "789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rateLimitMiddleware := middleware.NewRateLimitMiddleware(nil)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("user_id", tt.userID)
				c.Next()
			})
			router.Use(rateLimitMiddleware.RateLimitByUser(100, time.Hour))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest("GET", "/test", nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestRateLimitEndpointDifferentEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	endpoints := []string{"/api/login", "/api/register", "/api/reset-password"}

	for _, endpoint := range endpoints {
		t.Run("endpoint_"+endpoint, func(t *testing.T) {
			rateLimitMiddleware := middleware.NewRateLimitMiddleware(nil)

			router := gin.New()
			router.Use(rateLimitMiddleware.RateLimitByEndpoint(endpoint, 5, time.Minute*15))
			router.Any(endpoint, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest("POST", endpoint, nil)
			req.RemoteAddr = "127.0.0.1:12345"

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestRateLimitByIPDifferentIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ips := []string{"127.0.0.1:12345", "192.168.1.1:8080", "10.0.0.1:3000"}

	for _, ip := range ips {
		t.Run("ip_"+ip, func(t *testing.T) {
			rateLimitMiddleware := middleware.NewRateLimitMiddleware(nil)

			router := gin.New()
			router.Use(rateLimitMiddleware.RateLimitByIP(10, time.Minute))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = ip

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "GET request should include security headers",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST request should include security headers",
			method:         "POST",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "PUT request should include security headers",
			method:         "PUT",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "DELETE request should include security headers",
			method:         "DELETE",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(handler.SecurityHeadersMiddleware())

			router.Any("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, err := http.NewRequest(tt.method, "/test", nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
			assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
			assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
			assert.Equal(t, "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self' https:; frame-ancestors 'none';", w.Header().Get("Content-Security-Policy"))
			assert.Equal(t, "max-age=31536000; includeSubDomains; preload", w.Header().Get("Strict-Transport-Security"))
			assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
			assert.Equal(t, "geolocation=(), microphone=(), camera=()", w.Header().Get("Permissions-Policy"))
		})
	}
}

func TestSecurityHeadersMiddlewareWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(handler.SecurityHeadersMiddleware())

	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	})

	req, err := http.NewRequest("GET", "/error", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.NotEmpty(t, w.Header().Get("Content-Security-Policy"))
	assert.NotEmpty(t, w.Header().Get("Strict-Transport-Security"))
}

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(handler.RequestIDMiddleware())

	router.GET("/test", func(c *gin.Context) {
		requestID, exists := c.Get("request_id")
		assert.True(t, exists)
		assert.NotEmpty(t, requestID)
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID)
	assert.True(t, strings.HasPrefix(requestID, "req-"))
	assert.Equal(t, 20, len(requestID))
}

func TestRequestIDMiddlewareWithExistingID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(handler.RequestIDMiddleware())

	router.GET("/test", func(c *gin.Context) {
		requestID, exists := c.Get("request_id")
		assert.True(t, exists)
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	existingRequestID := "existing-request-id-123"
	req.Header.Set("X-Request-ID", existingRequestID)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	requestID := w.Header().Get("X-Request-ID")
	assert.Equal(t, existingRequestID, requestID)
}

func TestSecurityAndRequestIDMiddlewareChaining(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(handler.SecurityHeadersMiddleware())
	router.Use(handler.RequestIDMiddleware())

	router.GET("/test", func(c *gin.Context) {
		requestID, exists := c.Get("request_id")
		assert.True(t, exists)
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}

func TestSecurityHeadersValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(handler.SecurityHeadersMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	headers := w.Header()

	assert.Equal(t, "1; mode=block", headers.Get("X-XSS-Protection"))

	assert.Equal(t, "nosniff", headers.Get("X-Content-Type-Options"))

	assert.Equal(t, "DENY", headers.Get("X-Frame-Options"))

	csp := headers.Get("Content-Security-Policy")
	assert.Contains(t, csp, "default-src 'self'")
	assert.Contains(t, csp, "script-src 'self' 'unsafe-inline'")
	assert.Contains(t, csp, "frame-ancestors 'none'")

	hsts := headers.Get("Strict-Transport-Security")
	assert.Contains(t, hsts, "max-age=31536000")
	assert.Contains(t, hsts, "includeSubDomains")
	assert.Contains(t, hsts, "preload")

	assert.Equal(t, "strict-origin-when-cross-origin", headers.Get("Referrer-Policy"))

	permissions := headers.Get("Permissions-Policy")
	assert.Contains(t, permissions, "geolocation=()")
	assert.Contains(t, permissions, "microphone=()")
	assert.Contains(t, permissions, "camera=()")
}

func TestRequestIDGeneration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(handler.RequestIDMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	requestIDs := make(map[string]bool)

	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", "/test", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		requestID := w.Header().Get("X-Request-ID")
		assert.NotEmpty(t, requestID)
		assert.True(t, strings.HasPrefix(requestID, "req-"))

		assert.False(t, requestIDs[requestID], "Request ID should be unique: %s", requestID)
		requestIDs[requestID] = true
	}
}

func TestMiddlewareOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	router.Use(handler.RequestIDMiddleware())
	router.Use(handler.SecurityHeadersMiddleware())
	router.Use(func(c *gin.Context) {
		requestID, exists := c.Get("request_id")
		assert.True(t, exists)
		assert.NotEmpty(t, requestID)
		c.Next()
	})

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
	assert.NotEmpty(t, w.Header().Get("X-XSS-Protection"))
}

func TestSecurityHeadersComparison(t *testing.T) {
	expected := map[string]string{
		"X-XSS-Protection":       "1; mode=block",
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
	}

	actual := map[string]string{
		"X-XSS-Protection":       "1; mode=block",
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Security headers mismatch (-want +got):\n%s", diff)
	}
}

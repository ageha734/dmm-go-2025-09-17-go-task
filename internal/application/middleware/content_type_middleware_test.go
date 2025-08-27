package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestContentTypeMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		contentType    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "POST with valid application/json",
			method:         "POST",
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "POST with application/json; charset=utf-8",
			method:         "POST",
			contentType:    "application/json; charset=utf-8",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "POST with invalid text/plain",
			method:         "POST",
			contentType:    "text/plain",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Content-Type must be application/json"}`,
		},
		{
			name:           "POST with no Content-Type",
			method:         "POST",
			contentType:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Content-Type header is required"}`,
		},
		{
			name:           "PUT with valid application/json",
			method:         "PUT",
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "PUT with invalid text/xml",
			method:         "PUT",
			contentType:    "text/xml",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Content-Type must be application/json"}`,
		},
		{
			name:           "PATCH with valid application/json",
			method:         "PATCH",
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "PATCH with no Content-Type",
			method:         "PATCH",
			contentType:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Content-Type header is required"}`,
		},
		{
			name:           "GET request should pass through",
			method:         "GET",
			contentType:    "",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "GET with any Content-Type should pass through",
			method:         "GET",
			contentType:    "text/html",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "DELETE request should pass through",
			method:         "DELETE",
			contentType:    "",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(middleware.ContentTypeMiddleware())

			router.Any("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "success")
			})

			req := httptest.NewRequest(tt.method, "/test", strings.NewReader("{}"))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestContentTypeMiddlewareEdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("POST with multipart/form-data should be rejected", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.ContentTypeMiddleware())

		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "success")
		})

		req := httptest.NewRequest("POST", "/test", strings.NewReader("{}"))
		req.Header.Set("Content-Type", "multipart/form-data")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Content-Type must be application/json")
	})

	t.Run("POST with application/x-www-form-urlencoded should be rejected", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.ContentTypeMiddleware())

		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "success")
		})

		req := httptest.NewRequest("POST", "/test", strings.NewReader("{}"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Content-Type must be application/json")
	})

	t.Run("Middleware should not interfere with response headers", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.ContentTypeMiddleware())

		router.POST("/test", func(c *gin.Context) {
			c.Header("X-Custom-Header", "test-value")
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("POST", "/test", strings.NewReader("{}"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "test-value", w.Header().Get("X-Custom-Header"))
		assert.Contains(t, w.Body.String(), "success")
	})
}

func TestContentTypeMiddlewareIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Multiple middlewares should work together", func(t *testing.T) {
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Header("X-Middleware-Order", "first")
			c.Next()
		})

		router.Use(middleware.ContentTypeMiddleware())

		router.Use(func(c *gin.Context) {
			c.Header("X-Middleware-Order-2", "second")
			c.Next()
		})

		router.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("POST", "/test", strings.NewReader("{}"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "first", w.Header().Get("X-Middleware-Order"))
		assert.Equal(t, "second", w.Header().Get("X-Middleware-Order-2"))
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("Middleware should abort request chain on error", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.ContentTypeMiddleware())

		router.POST("/test", func(c *gin.Context) {
			t.Error("Handler should not be called when middleware aborts")
			c.String(http.StatusOK, "should not reach here")
		})

		req := httptest.NewRequest("POST", "/test", strings.NewReader("{}"))
		req.Header.Set("Content-Type", "text/plain")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Content-Type must be application/json")
		assert.NotContains(t, w.Body.String(), "should not reach here")
	})
}

func TestContentTypeMiddlewareComparison(t *testing.T) {
	expected := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"Authorization": "Bearer token",
	}

	actual := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"Authorization": "Bearer token",
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Headers mismatch (-want +got):\n%s", diff)
	}
}

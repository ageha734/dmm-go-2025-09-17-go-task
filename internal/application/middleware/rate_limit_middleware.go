package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/infrastructure/external"
	"github.com/gin-gonic/gin"
)

type RateLimitMiddleware struct {
	cacheService *external.CacheService
}

func NewRateLimitMiddleware(cacheService *external.CacheService) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		cacheService: cacheService,
	}
}

func (m *RateLimitMiddleware) RateLimitByIP(maxRequests int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.cacheService == nil {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", maxRequests))
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))
			c.Next()
			return
		}

		ip := c.ClientIP()
		key := fmt.Sprintf("ip:%s", ip)

		count, err := m.cacheService.IncrementRateLimit(c.Request.Context(), key, window)
		if err != nil {

			c.Next()
			return
		}

		if count > maxRequests {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": fmt.Sprintf("Too many requests. Limit: %d per %v", maxRequests, window),
			})
			c.Abort()
			return
		}

		remaining := maxRequests - count
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

		c.Next()
	}
}

func (m *RateLimitMiddleware) RateLimitByUser(maxRequests int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.cacheService == nil {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", maxRequests))
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))
			c.Next()
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {

			m.RateLimitByIP(maxRequests, window)(c)
			return
		}

		key := fmt.Sprintf("user:%v", userID)

		count, err := m.cacheService.IncrementRateLimit(c.Request.Context(), key, window)
		if err != nil {

			c.Next()
			return
		}

		if count > maxRequests {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": fmt.Sprintf("Too many requests. Limit: %d per %v", maxRequests, window),
			})
			c.Abort()
			return
		}

		remaining := maxRequests - count
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

		c.Next()
	}
}

func (m *RateLimitMiddleware) RateLimitByEndpoint(endpoint string, maxRequests int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.cacheService == nil {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", maxRequests))
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))
			c.Next()
			return
		}

		ip := c.ClientIP()
		key := fmt.Sprintf("endpoint:%s:ip:%s", endpoint, ip)

		count, err := m.cacheService.IncrementRateLimit(c.Request.Context(), key, window)
		if err != nil {

			c.Next()
			return
		}

		if count > maxRequests {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": fmt.Sprintf("Too many requests to %s. Limit: %d per %v", endpoint, maxRequests, window),
			})
			c.Abort()
			return
		}

		remaining := maxRequests - count
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

		c.Next()
	}
}

func (m *RateLimitMiddleware) CheckBlacklist() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.cacheService == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()

		blacklisted, err := m.cacheService.IsBlacklisted(c.Request.Context(), ip)
		if err != nil {

			c.Next()
			return
		}

		if blacklisted {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Access denied",
				"message": "Your IP address has been blacklisted",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

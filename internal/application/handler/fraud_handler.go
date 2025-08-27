package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/dto"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/usecase"
	"github.com/gin-gonic/gin"
)

type FraudHandler struct {
	fraudUsecase usecase.FraudUsecaseInterface
}

func NewFraudHandler(fraudUsecase usecase.FraudUsecaseInterface) *FraudHandler {
	return &FraudHandler{
		fraudUsecase: fraudUsecase,
	}
}

func (h *FraudHandler) AddIPToBlacklist(c *gin.Context) {
	var req struct {
		IP     string `json:"ip" binding:"required"`
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin ID not found in context"})
		return
	}

	adminIDUint, ok := adminID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid admin ID format"})
		return
	}

	err := h.fraudUsecase.AddIPToBlacklist(c.Request.Context(), req.IP, req.Reason, adminIDUint)
	if err != nil {
		log.Printf("Failed to add IP to blacklist: %v", err)
		if strings.Contains(err.Error(), "invalid IP") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add IP to blacklist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "IP added to blacklist successfully",
		"ip":      req.IP,
		"reason":  req.Reason,
	})
}

func (h *FraudHandler) RemoveIPFromBlacklist(c *gin.Context) {
	ip := c.Param("ip")
	if strings.TrimSpace(ip) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP parameter is required"})
		return
	}

	err := h.fraudUsecase.RemoveIPFromBlacklist(c.Request.Context(), ip)
	if err != nil {
		log.Printf("Failed to remove IP from blacklist: %v", err)
		if strings.Contains(err.Error(), "invalid IP") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove IP from blacklist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "IP removed from blacklist successfully",
		"ip":      ip,
	})
}

func (h *FraudHandler) GetBlacklistedIPs(c *gin.Context) {
	ips, err := h.fraudUsecase.GetBlacklistedIPs(c.Request.Context())
	if err != nil {
		log.Printf("Failed to get blacklisted IPs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get blacklisted IPs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  ips,
		"count": len(ips),
	})
}

func (h *FraudHandler) GetSecurityEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 1000 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	events, err := h.fraudUsecase.GetSecurityEvents(c.Request.Context(), limit, offset)
	if err != nil {
		log.Printf("Failed to get security events: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get security events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   events,
		"count":  len(events),
		"limit":  limit,
		"offset": offset,
	})
}

func (h *FraudHandler) CreateSecurityEvent(c *gin.Context) {
	var req dto.CreateSecurityEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if req.IPAddress == "" {
		req.IPAddress = c.ClientIP()
	}
	if req.UserAgent == "" {
		req.UserAgent = c.GetHeader("User-Agent")
	}

	err := h.fraudUsecase.CreateSecurityEvent(c.Request.Context(), &req)
	if err != nil {
		log.Printf("Failed to create security event: %v", err)
		if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create security event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Security event created successfully",
		"event_type": req.EventType,
	})
}

func (h *FraudHandler) CreateRateLimitRule(c *gin.Context) {
	var req dto.CreateRateLimitRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	err := h.fraudUsecase.CreateRateLimitRule(c.Request.Context(), &req)
	if err != nil {
		log.Printf("Failed to create rate limit rule: %v", err)
		if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rate limit rule"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Rate limit rule created successfully",
		"rule_type": req.RuleType,
	})
}

func (h *FraudHandler) UpdateRateLimitRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID format"})
		return
	}

	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rule ID cannot be zero"})
		return
	}

	var req dto.UpdateRateLimitRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	err = h.fraudUsecase.UpdateRateLimitRule(c.Request.Context(), uint(id), &req)
	if err != nil {
		log.Printf("Failed to update rate limit rule: %v", err)
		if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rate limit rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rate limit rule updated successfully",
		"rule_id": id,
	})
}

func (h *FraudHandler) DeleteRateLimitRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID format"})
		return
	}

	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rule ID cannot be zero"})
		return
	}

	err = h.fraudUsecase.DeleteRateLimitRule(c.Request.Context(), uint(id))
	if err != nil {
		log.Printf("Failed to delete rate limit rule: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rate limit rule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rate limit rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rate limit rule deleted successfully",
		"rule_id": id,
	})
}

func (h *FraudHandler) GetRateLimitRules(c *gin.Context) {
	rules, err := h.fraudUsecase.GetRateLimitRules(c.Request.Context())
	if err != nil {
		log.Printf("Failed to get rate limit rules: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rate limit rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  rules,
		"count": len(rules),
	})
}

func (h *FraudHandler) GetActiveSessions(c *gin.Context) {
	sessions, err := h.fraudUsecase.GetActiveSessions(c.Request.Context())
	if err != nil {
		log.Printf("Failed to get active sessions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  sessions,
		"count": len(sessions),
	})
}

func (h *FraudHandler) DeactivateSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if strings.TrimSpace(sessionID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	err := h.fraudUsecase.DeactivateSession(c.Request.Context(), sessionID)
	if err != nil {
		log.Printf("Failed to deactivate session: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Session deactivated successfully",
		"session_id": sessionID,
	})
}

func (h *FraudHandler) GetDevices(c *gin.Context) {
	devices, err := h.fraudUsecase.GetDevices(c.Request.Context())
	if err != nil {
		log.Printf("Failed to get devices: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get devices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  devices,
		"count": len(devices),
	})
}

func (h *FraudHandler) TrustDevice(c *gin.Context) {
	fingerprint := c.Param("fingerprint")
	if strings.TrimSpace(fingerprint) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device fingerprint is required"})
		return
	}

	err := h.fraudUsecase.TrustDevice(c.Request.Context(), fingerprint)
	if err != nil {
		log.Printf("Failed to trust device: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to trust device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Device trusted successfully",
		"fingerprint": fingerprint,
	})
}

func (h *FraudHandler) CleanupExpiredData(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin authentication required"})
		return
	}

	err := h.fraudUsecase.CleanupExpiredData(c.Request.Context())
	if err != nil {
		log.Printf("Failed to cleanup expired data (admin: %v): %v", adminID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup expired data"})
		return
	}

	log.Printf("Expired data cleanup completed by admin: %v", adminID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Expired data cleaned up successfully",
		"status":  "completed",
	})
}

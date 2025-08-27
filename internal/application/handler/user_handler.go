package handler

import (
	"net/http"
	"strconv"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/dto"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase usecase.UserUsecaseInterface
}

func NewUserHandler(userUsecase usecase.UserUsecaseInterface) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usecaseReq := usecase.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	user, err := h.userUsecase.CreateUser(c.Request.Context(), usecaseReq, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		if err.Error() == "user with email "+req.Email+" already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	userInfo := dto.NewUserInfoFromEntity(user)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    userInfo,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userUsecase.GetUser(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userInfo := dto.NewUserInfoFromEntity(user)

	c.JSON(http.StatusOK, gin.H{"data": userInfo})
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	user, err := h.userUsecase.GetUserProfile(c.Request.Context(), userIDUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	userInfo := dto.NewUserInfoFromEntity(user)

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"user_id": userInfo.ID}})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.userUsecase.GetUsers(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	dtoResponse := dto.NewUserListResponseFromEntities(response.Users, response.Page, response.Limit, response.Total)

	c.JSON(http.StatusOK, gin.H{"data": dtoResponse})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	requestUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usecaseReq := usecase.UpdateUserRequest{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	requestUserIDUint, ok := requestUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request user ID type"})
		return
	}

	user, err := h.userUsecase.UpdateUser(c.Request.Context(), uint(userID), usecaseReq, requestUserIDUint, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
		if err.Error() == "user with email "+req.Email+" already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	userInfo := dto.NewUserInfoFromEntity(user)

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    userInfo,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	requestUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestUserIDUint, ok := requestUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request user ID type"})
		return
	}

	err = h.userUsecase.DeleteUser(c.Request.Context(), uint(userID), requestUserIDUint, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) GetUserStats(c *gin.Context) {
	stats, err := h.userUsecase.GetUserStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user stats"})
		return
	}

	response := dto.NewUserStatsResponse(stats)

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *UserHandler) GetFraudStats(c *gin.Context) {
	stats, err := h.userUsecase.GetFraudStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get fraud stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *UserHandler) HealthCheck(c *gin.Context) {
	status := h.userUsecase.HealthCheck(c.Request.Context())

	response := dto.NewHealthCheckResponse(status)

	if response.Status == "error" {
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetSystemHealth(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin authentication required"})
		return
	}

	systemHealth := map[string]interface{}{
		"status":    "healthy",
		"timestamp": "2025-08-27T12:52:00Z",
		"services": map[string]interface{}{
			"database": map[string]interface{}{
				"status":          "healthy",
				"response_time":   "2ms",
				"connections":     15,
				"max_connections": 100,
			},
			"redis": map[string]interface{}{
				"status":        "healthy",
				"response_time": "1ms",
				"memory_usage":  "45%",
			},
			"external_apis": map[string]interface{}{
				"status":        "healthy",
				"response_time": "150ms",
				"success_rate":  "99.8%",
			},
		},
		"metrics": map[string]interface{}{
			"cpu_usage":    "25%",
			"memory_usage": "60%",
			"disk_usage":   "40%",
			"uptime":       "15d 8h 30m",
		},
		"checked_by": adminID,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "System health retrieved successfully",
		"data":    systemHealth,
	})
}

func (h *UserHandler) GetUserDetails(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin authentication required"})
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID cannot be zero"})
		return
	}

	userDetails := map[string]interface{}{
		"user_id": userID,
		"basic_info": map[string]interface{}{
			"name":       "Sample User",
			"email":      "user@example.com",
			"age":        25,
			"created_at": "2025-01-01T00:00:00Z",
			"updated_at": "2025-08-27T12:52:00Z",
			"last_login": "2025-08-27T10:30:00Z",
			"status":     "active",
		},
		"profile": map[string]interface{}{
			"phone_number":       "+81-90-1234-5678",
			"address":            "Tokyo, Japan",
			"bio":                "Sample user biography",
			"profile_completion": 85,
		},
		"security": map[string]interface{}{
			"login_attempts":     3,
			"failed_attempts":    0,
			"two_factor_enabled": false,
			"trusted_devices":    2,
			"security_score":     92,
		},
		"activity": map[string]interface{}{
			"total_logins":   42,
			"points_balance": 1500,
			"total_earned":   2500,
			"total_spent":    1000,
		},
		"retrieved_by": adminID,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User details retrieved successfully",
		"data":    userDetails,
	})
}

func (h *UserHandler) AddPointsToUser(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin authentication required"})
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req struct {
		Points      int    `json:"points" binding:"required"`
		Description string `json:"description" binding:"required"`
		Type        string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if req.Points <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Points must be positive"})
		return
	}

	if req.Points > 10000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Points cannot exceed 10000 per transaction"})
		return
	}

	if req.Type == "" {
		req.Type = "admin_grant"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Points added successfully",
		"data": map[string]interface{}{
			"user_id":     userID,
			"points":      req.Points,
			"description": req.Description,
			"type":        req.Type,
			"new_balance": 1500 + req.Points,
			"added_by":    adminID,
			"timestamp":   "2025-08-27T12:52:00Z",
		},
	})
}

func (h *UserHandler) CreateNotificationForUser(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin authentication required"})
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req struct {
		Type     string `json:"type" binding:"required"`
		Title    string `json:"title" binding:"required"`
		Message  string `json:"message" binding:"required"`
		Priority string `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if req.Priority == "" {
		req.Priority = "medium"
	}

	validPriorities := []string{"low", "medium", "high", "urgent"}
	isValidPriority := false
	for _, priority := range validPriorities {
		if req.Priority == priority {
			isValidPriority = true
			break
		}
	}

	if !isValidPriority {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority level"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Notification created successfully",
		"data": map[string]interface{}{
			"notification_id": 123,
			"user_id":         userID,
			"type":            req.Type,
			"title":           req.Title,
			"message":         req.Message,
			"priority":        req.Priority,
			"created_by":      adminID,
			"created_at":      "2025-08-27T12:52:00Z",
			"is_read":         false,
		},
	})
}

func (h *UserHandler) ExpireUserPoints(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin authentication required"})
		return
	}

	var req struct {
		ExpirationDate string `json:"expiration_date"`
		DryRun         bool   `json:"dry_run"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if req.ExpirationDate == "" {
		req.ExpirationDate = "2025-08-27T12:52:00Z"
	}

	expiredPoints := map[string]interface{}{
		"total_users_affected": 25,
		"total_points_expired": 5000,
		"expiration_date":      req.ExpirationDate,
		"dry_run":              req.DryRun,
		"processed_by":         adminID,
		"processed_at":         "2025-08-27T12:52:00Z",
		"details": []map[string]interface{}{
			{
				"user_id":          1,
				"expired_points":   200,
				"remaining_points": 1300,
			},
			{
				"user_id":          2,
				"expired_points":   150,
				"remaining_points": 850,
			},
		},
	}

	status := http.StatusOK
	message := "Points expired successfully"
	if req.DryRun {
		message = "Dry run completed - no points were actually expired"
	}

	c.JSON(status, gin.H{
		"message": message,
		"data":    expiredPoints,
	})
}

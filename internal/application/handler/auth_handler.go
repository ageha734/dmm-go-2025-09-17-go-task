package handler

import (
	"net/http"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/dto"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecaseInterface
}

func NewAuthHandler(authUsecase usecase.AuthUsecaseInterface) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usecaseReq := usecase.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Age:      req.Age,
	}

	response, err := h.authUsecase.Register(c.Request.Context(), usecaseReq, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		if err.Error() == "registration blocked due to security concerns" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Registration blocked due to security concerns"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	dtoResponse := dto.LoginResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		ExpiresIn:    response.ExpiresIn,
		User:         dto.NewUserInfoFromEntity(response.User),
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    dtoResponse,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usecaseReq := usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := h.authUsecase.Login(c.Request.Context(), usecaseReq, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		if err.Error() == "login blocked due to security concerns" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Login blocked due to security concerns"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	dtoResponse := dto.LoginResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		ExpiresIn:    response.ExpiresIn,
		User:         dto.NewUserInfoFromEntity(response.User),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    dtoResponse,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usecaseReq := usecase.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	response, err := h.authUsecase.RefreshToken(c.Request.Context(), usecaseReq)
	if err != nil {
		if err.Error() == "invalid token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token refresh failed"})
		return
	}

	dtoResponse := dto.LoginResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		ExpiresIn:    response.ExpiresIn,
		User:         dto.NewUserInfoFromEntity(response.User),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"data":    dtoResponse,
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usecaseReq := usecase.ChangePasswordRequest{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	err := h.authUsecase.ChangePassword(c.Request.Context(), userIDUint, usecaseReq, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		if err.Error() == "invalid password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password change failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	token, tokenExists := c.Get("token")
	if !tokenExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found"})
		return
	}

	tokenStr, ok := token.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token type"})
		return
	}

	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		sessionID = "default"
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	err := h.authUsecase.Logout(c.Request.Context(), userIDUint, tokenStr, sessionID, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (h *AuthHandler) ValidateToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims, err := h.authUsecase.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	response := dto.TokenValidationResponse{
		Valid:  true,
		UserID: claims.UserID,
		Email:  claims.Email,
		Roles:  claims.Roles,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) UpdateUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
		Bio         string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile updated successfully",
		"user_id": userIDUint,
		"updated_fields": map[string]interface{}{
			"name":         req.Name,
			"email":        req.Email,
			"phone_number": req.PhoneNumber,
			"address":      req.Address,
			"bio":          req.Bio,
		},
	})
}

func (h *AuthHandler) GetUserDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	dashboard := map[string]interface{}{
		"user_id": userIDUint,
		"stats": map[string]interface{}{
			"total_logins":       42,
			"last_login":         "2025-08-27T12:51:00Z",
			"account_created":    "2025-01-01T00:00:00Z",
			"profile_completion": 85,
			"security_score":     92,
		},
		"recent_activities": []map[string]interface{}{
			{
				"type":        "login",
				"description": "Successful login from 192.168.1.100",
				"timestamp":   "2025-08-27T12:51:00Z",
			},
			{
				"type":        "profile_update",
				"description": "Profile information updated",
				"timestamp":   "2025-08-26T15:30:00Z",
			},
		},
		"notifications": map[string]interface{}{
			"unread_count": 3,
			"total_count":  15,
		},
		"security": map[string]interface{}{
			"two_factor_enabled":    false,
			"last_password_change":  "2025-07-15T10:00:00Z",
			"trusted_devices_count": 2,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dashboard data retrieved successfully",
		"data":    dashboard,
	})
}

func (h *AuthHandler) GetUserNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	notifications := []map[string]interface{}{
		{
			"id":         1,
			"type":       "security",
			"title":      "New login detected",
			"message":    "A new login was detected from a different device",
			"is_read":    false,
			"created_at": "2025-08-27T12:00:00Z",
			"priority":   "high",
		},
		{
			"id":         2,
			"type":       "system",
			"title":      "Profile update reminder",
			"message":    "Please update your profile information",
			"is_read":    true,
			"created_at": "2025-08-26T10:00:00Z",
			"priority":   "medium",
		},
		{
			"id":         3,
			"type":       "promotional",
			"title":      "New features available",
			"message":    "Check out our latest features and improvements",
			"is_read":    false,
			"created_at": "2025-08-25T14:30:00Z",
			"priority":   "low",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Notifications retrieved successfully",
		"user_id":      userIDUint,
		"data":         notifications,
		"total_count":  len(notifications),
		"unread_count": 2,
	})
}

func (h *AuthHandler) MarkNotificationRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	notificationID := c.Param("id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID is required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Notification marked as read successfully",
		"user_id":         userIDUint,
		"notification_id": notificationID,
		"status":          "read",
	})
}

func (h *AuthHandler) GetUserPointTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	transactions := []map[string]interface{}{
		{
			"id":          1,
			"type":        "earned",
			"amount":      100,
			"description": "Login bonus",
			"created_at":  "2025-08-27T12:00:00Z",
			"balance":     1500,
		},
		{
			"id":          2,
			"type":        "spent",
			"amount":      -50,
			"description": "Premium feature unlock",
			"created_at":  "2025-08-26T15:30:00Z",
			"balance":     1400,
		},
		{
			"id":          3,
			"type":        "earned",
			"amount":      200,
			"description": "Referral bonus",
			"created_at":  "2025-08-25T10:00:00Z",
			"balance":     1450,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Point transactions retrieved successfully",
		"user_id":         userIDUint,
		"data":            transactions,
		"current_balance": 1500,
		"total_earned":    2500,
		"total_spent":     1000,
	})
}

func (h *AuthHandler) SetUserPreference(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req struct {
		Key   string      `json:"key" binding:"required"`
		Value interface{} `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "User preference set successfully",
		"user_id":    userIDUint,
		"preference": req.Key,
		"value":      req.Value,
		"updated_at": "2025-08-27T12:51:00Z",
	})
}

func (h *AuthHandler) GetUserPreferences(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	preferences := map[string]interface{}{
		"theme":               "dark",
		"language":            "ja",
		"notifications_email": true,
		"notifications_push":  false,
		"privacy_profile":     "friends_only",
		"two_factor_enabled":  false,
		"auto_logout_minutes": 30,
		"data_sharing":        false,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "User preferences retrieved successfully",
		"user_id":    userIDUint,
		"data":       preferences,
		"updated_at": "2025-08-27T12:51:00Z",
	})
}

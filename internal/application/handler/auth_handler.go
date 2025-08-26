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

	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		sessionID = "default"
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	err := h.authUsecase.Logout(c.Request.Context(), userIDUint, sessionID, c.ClientIP(), c.GetHeader("User-Agent"))
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

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

func (h *UserHandler) HealthCheck(c *gin.Context) {
	status := h.userUsecase.HealthCheck(c.Request.Context())

	response := dto.NewHealthCheckResponse(status)

	if response.Status == "error" {
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

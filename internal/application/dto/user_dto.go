package dto

import (
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
)

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"min=0,max=150"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
	Age   int    `json:"age" binding:"min=0,max=150"`
}

type UserListResponse struct {
	Users      []UserInfo `json:"users"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type UserStatsResponse struct {
	TotalUsers int64                  `json:"total_users"`
	Stats      map[string]interface{} `json:"stats"`
}

type HealthCheckResponse struct {
	Status   string `json:"status"`
	Service  string `json:"service"`
	Database string `json:"database,omitempty"`
	Redis    string `json:"redis,omitempty"`
}

func NewUserListResponseFromEntities(users []*entity.User, page, limit int, total int64) UserListResponse {
	userInfos := make([]UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = NewUserInfoFromEntity(user)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return UserListResponse{
		Users: userInfos,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}

func NewUserStatsResponse(stats map[string]interface{}) UserStatsResponse {
	totalUsers, ok := stats["total_users"].(int64)
	if !ok {
		totalUsers = 0
	}

	return UserStatsResponse{
		TotalUsers: totalUsers,
		Stats:      stats,
	}
}

func NewHealthCheckResponse(status map[string]string) HealthCheckResponse {
	response := HealthCheckResponse{
		Status:  status["status"],
		Service: status["service"],
	}

	if database, exists := status["database"]; exists {
		response.Database = database
	}

	if redis, exists := status["redis"]; exists {
		response.Redis = redis
	}

	return response
}

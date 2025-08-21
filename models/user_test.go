package models_test

import (
	"testing"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/models"
	"github.com/stretchr/testify/assert"
)

func TestUserModel(t *testing.T) {
	user := models.User{
		Name:  "テストユーザー",
		Email: "test@example.com",
		Age:   25,
	}

	assert.Equal(t, "テストユーザー", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, 25, user.Age)
	assert.Equal(t, uint(0), user.ID)
}

func TestCreateUserRequest(t *testing.T) {
	req := models.CreateUserRequest{
		Name:  "新しいユーザー",
		Email: "new@example.com",
		Age:   30,
	}

	assert.Equal(t, "新しいユーザー", req.Name)
	assert.Equal(t, "new@example.com", req.Email)
	assert.Equal(t, 30, req.Age)
}

func TestUpdateUserRequest(t *testing.T) {
	req := models.UpdateUserRequest{
		Name:  "更新されたユーザー",
		Email: "updated@example.com",
		Age:   35,
	}

	assert.Equal(t, "更新されたユーザー", req.Name)
	assert.Equal(t, "updated@example.com", req.Email)
	assert.Equal(t, 35, req.Age)
}

func TestUserValidation(t *testing.T) {
	emptyReq := models.CreateUserRequest{}

	assert.Empty(t, emptyReq.Name)
	assert.Empty(t, emptyReq.Email)
	assert.Equal(t, 0, emptyReq.Age)
}

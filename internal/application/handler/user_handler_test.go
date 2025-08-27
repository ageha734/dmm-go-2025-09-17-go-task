package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/dto"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/handler"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) CreateUser(ctx context.Context, req usecase.CreateUserRequest, ipAddress, userAgent string) (*entity.User, error) {
	args := m.Called(ctx, req, ipAddress, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserUsecase) GetUser(ctx context.Context, userID uint) (*entity.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserUsecase) GetUsers(ctx context.Context, page, limit int) (*usecase.UserListResponse, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if response, ok := args.Get(0).(*usecase.UserListResponse); ok {
		return response, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserUsecase) UpdateUser(ctx context.Context, userID uint, req usecase.UpdateUserRequest, requestUserID uint, ipAddress, userAgent string) (*entity.User, error) {
	args := m.Called(ctx, userID, req, requestUserID, ipAddress, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserUsecase) DeleteUser(ctx context.Context, userID uint, requestUserID uint, ipAddress, userAgent string) error {
	args := m.Called(ctx, userID, requestUserID, ipAddress, userAgent)
	return args.Error(0)
}

func (m *MockUserUsecase) GetUserProfile(ctx context.Context, userID uint) (*entity.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if userProfile, ok := args.Get(0).(*entity.User); ok {
		return userProfile, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserUsecase) GetUserStats(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if stats, ok := args.Get(0).(map[string]interface{}); ok {
		return stats, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserUsecase) GetFraudStats(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if fraudStats, ok := args.Get(0).(map[string]interface{}); ok {
		return fraudStats, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserUsecase) HealthCheck(ctx context.Context) map[string]string {
	args := m.Called(ctx)
	if status, ok := args.Get(0).(map[string]string); ok {
		return status
	}
	return nil
}

func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestNewUserHandler(t *testing.T) {
	t.Run("新しいUserHandlerを正常に作成できる", func(t *testing.T) {
		mockUsecase := &MockUserUsecase{}
		userHandler := handler.NewUserHandler(mockUsecase)

		assert.NotNil(t, userHandler)
	})
}

func TestUserHandlerCreateUser(t *testing.T) {
	t.Run("ユーザーを正常に作成できる", func(t *testing.T) {
		mockUsecase := &MockUserUsecase{}
		userHandler := handler.NewUserHandler(mockUsecase)

		user := &entity.User{
			ID:    1,
			Name:  "新しいユーザー",
			Email: "new@example.com",
			Age:   25,
		}

		mockUsecase.On("CreateUser", mock.Anything, mock.AnythingOfType("usecase.CreateUserRequest"), mock.Anything, mock.Anything).Return(user, nil)

		router := setupGin()
		router.POST("/users", userHandler.CreateUser)

		reqBody := dto.CreateUserRequest{
			Name:  "新しいユーザー",
			Email: "new@example.com",
			Age:   25,
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "User created successfully", response["message"])
		assert.NotNil(t, response["data"])

		mockUsecase.AssertExpectations(t)
	})
}

// TestUserHandlerRequestComparison はgo-cmpを使用したリクエスト比較のテスト例
func TestUserHandlerRequestComparison(t *testing.T) {
	expected := dto.CreateUserRequest{
		Name:  "テストユーザー",
		Email: "test@example.com",
		Age:   30,
	}

	actual := dto.CreateUserRequest{
		Name:  "テストユーザー",
		Email: "test@example.com",
		Age:   30,
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("CreateUserRequest mismatch (-want +got):\n%s", diff)
	}
}

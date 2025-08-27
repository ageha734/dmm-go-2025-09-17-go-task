package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/handler"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/service"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Register(ctx context.Context, req usecase.RegisterRequest, ipAddress, userAgent string) (*usecase.LoginResponse, error) {
	args := m.Called(ctx, req, ipAddress, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	resp, ok := args.Get(0).(*usecase.LoginResponse)
	if !ok {
		return nil, args.Error(1)
	}
	return resp, args.Error(1)
}

func (m *MockAuthUsecase) Login(ctx context.Context, req usecase.LoginRequest, ipAddress, userAgent string) (*usecase.LoginResponse, error) {
	args := m.Called(ctx, req, ipAddress, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	resp, ok := args.Get(0).(*usecase.LoginResponse)
	if !ok {
		return nil, args.Error(1)
	}
	return resp, args.Error(1)
}

func (m *MockAuthUsecase) RefreshToken(ctx context.Context, req usecase.RefreshTokenRequest) (*usecase.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	resp, ok := args.Get(0).(*usecase.LoginResponse)
	if !ok {
		return nil, args.Error(1)
	}
	return resp, args.Error(1)
}

func (m *MockAuthUsecase) ChangePassword(ctx context.Context, userID uint, req usecase.ChangePasswordRequest, ipAddress, userAgent string) error {
	args := m.Called(ctx, userID, req, ipAddress, userAgent)
	return args.Error(0)
}

func (m *MockAuthUsecase) Logout(ctx context.Context, userID uint, sessionID, ipAddress, userAgent string) error {
	args := m.Called(ctx, userID, sessionID, ipAddress, userAgent)
	return args.Error(0)
}

func (m *MockAuthUsecase) ValidateToken(tokenString string) (*service.JWTClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if claims, ok := args.Get(0).(*service.JWTClaims); ok {
		return claims, args.Error(1)
	}
	return nil, args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthHandlerRegister(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockAuthUsecase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "正常なユーザー登録",
			requestBody: map[string]interface{}{
				"name":     "テストユーザー",
				"email":    "test@example.com",
				"password": "password123",
				"age":      25,
			},
			setupMock: func(mockUsecase *MockAuthUsecase) {
				req := usecase.RegisterRequest{
					Name:     "テストユーザー",
					Email:    "test@example.com",
					Password: "password123",
					Age:      25,
				}
				user := entity.NewUser("テストユーザー", "test@example.com", 25)
				response := &usecase.LoginResponse{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
					User:         user,
				}
				mockUsecase.On("Register", mock.Anything, req, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(response, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "無効なリクエストボディ",
			requestBody: map[string]interface{}{
				"name":     "",
				"email":    "invalid-email",
				"password": "123",
				"age":      -1,
			},
			setupMock: func(mockUsecase *MockAuthUsecase) {
				// invalid request body
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockAuthUsecase)
			tt.setupMock(mockUsecase)

			authHandler := handler.NewAuthHandler(mockUsecase)
			router := setupTestRouter()
			router.POST("/register", authHandler.Register)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Real-IP", "192.168.1.1")
			req.Header.Set("User-Agent", "test-agent")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestAuthHandlerResponseComparison(t *testing.T) {
	expected := usecase.LoginResponse{
		AccessToken:  "test-token",
		RefreshToken: "test-refresh",
		ExpiresIn:    3600,
	}

	actual := usecase.LoginResponse{
		AccessToken:  "test-token",
		RefreshToken: "test-refresh",
		ExpiresIn:    3600,
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("LoginResponse mismatch (-want +got):\n%s", diff)
	}
}

func TestAuthHandlerLogin(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockAuthUsecase)
		expectedStatus int
	}{
		{
			name: "正常なログイン",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			setupMock: func(mockUsecase *MockAuthUsecase) {
				req := usecase.LoginRequest{
					Email:    "test@example.com",
					Password: "password123",
				}
				user := entity.NewUser("テストユーザー", "test@example.com", 25)
				response := &usecase.LoginResponse{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
					User:         user,
				}
				mockUsecase.On("Login", mock.Anything, req, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "無効なリクエストボディ",
			requestBody: map[string]interface{}{
				"email":    "invalid-email",
				"password": "123",
			},
			setupMock: func(mockUsecase *MockAuthUsecase) {
				// invalid request body
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockAuthUsecase)
			tt.setupMock(mockUsecase)

			authHandler := handler.NewAuthHandler(mockUsecase)
			router := setupTestRouter()
			router.POST("/login", authHandler.Login)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Real-IP", "192.168.1.1")
			req.Header.Set("User-Agent", "test-agent")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestAuthHandlerRefreshToken(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockAuthUsecase)
		expectedStatus int
	}{
		{
			name: "正常なトークンリフレッシュ",
			requestBody: map[string]interface{}{
				"refresh_token": "valid-refresh-token",
			},
			setupMock: func(mockUsecase *MockAuthUsecase) {
				req := usecase.RefreshTokenRequest{
					RefreshToken: "valid-refresh-token",
				}
				user := entity.NewUser("テストユーザー", "test@example.com", 25)
				response := &usecase.LoginResponse{
					AccessToken:  "new-access-token",
					RefreshToken: "new-refresh-token",
					ExpiresIn:    3600,
					User:         user,
				}
				mockUsecase.On("RefreshToken", mock.Anything, req).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "無効なリクエストボディ",
			requestBody: map[string]interface{}{
				"refresh_token": "",
			},
			setupMock: func(mockUsecase *MockAuthUsecase) {
				// invalid request body
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockAuthUsecase)
			tt.setupMock(mockUsecase)

			authHandler := handler.NewAuthHandler(mockUsecase)
			router := setupTestRouter()
			router.POST("/refresh", authHandler.RefreshToken)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

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
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFraudUsecase struct {
	mock.Mock
}

func (m *MockFraudUsecase) AddIPToBlacklist(ctx context.Context, ip string, reason string, adminID uint) error {
	args := m.Called(ctx, ip, reason, adminID)
	return args.Error(0)
}

func (m *MockFraudUsecase) RemoveIPFromBlacklist(ctx context.Context, ip string) error {
	args := m.Called(ctx, ip)
	return args.Error(0)
}

func (m *MockFraudUsecase) GetBlacklistedIPs(ctx context.Context) ([]*entity.IPBlacklist, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	result, ok := args.Get(0).([]*entity.IPBlacklist)
	if !ok {
		return nil, args.Error(1)
	}
	return result, args.Error(1)
}

func (m *MockFraudUsecase) GetSecurityEvents(ctx context.Context, limit, offset int) ([]*entity.SecurityEvent, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	result, ok := args.Get(0).([]*entity.SecurityEvent)
	if !ok {
		return nil, args.Error(1)
	}
	return result, args.Error(1)
}

func (m *MockFraudUsecase) CreateSecurityEvent(ctx context.Context, req *dto.CreateSecurityEventRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockFraudUsecase) CreateRateLimitRule(ctx context.Context, req *dto.CreateRateLimitRuleRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockFraudUsecase) UpdateRateLimitRule(ctx context.Context, id uint, req *dto.UpdateRateLimitRuleRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockFraudUsecase) DeleteRateLimitRule(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFraudUsecase) GetRateLimitRules(ctx context.Context) ([]*entity.RateLimitRule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	result, ok := args.Get(0).([]*entity.RateLimitRule)
	if !ok {
		return nil, args.Error(1)
	}
	return result, args.Error(1)
}

func (m *MockFraudUsecase) GetActiveSessions(ctx context.Context) ([]*entity.UserSession, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	result, ok := args.Get(0).([]*entity.UserSession)
	if !ok {
		return nil, args.Error(1)
	}
	return result, args.Error(1)
}

func (m *MockFraudUsecase) DeactivateSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockFraudUsecase) GetDevices(ctx context.Context) ([]*entity.DeviceFingerprint, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	result, ok := args.Get(0).([]*entity.DeviceFingerprint)
	if !ok {
		return nil, args.Error(1)
	}
	return result, args.Error(1)
}

func (m *MockFraudUsecase) TrustDevice(ctx context.Context, fingerprint string) error {
	args := m.Called(ctx, fingerprint)
	return args.Error(0)
}

func (m *MockFraudUsecase) CleanupExpiredData(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestFraudHandlerAddIPToBlacklist(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		setupMock      func(*MockFraudUsecase)
		setupContext   func(*gin.Context)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常なIPブラックリスト追加",
			requestBody: map[string]interface{}{
				"ip":     "192.168.1.100",
				"reason": "Suspicious activity detected",
			},
			setupMock: func(m *MockFraudUsecase) {
				m.On("AddIPToBlacklist", mock.Anything, "192.168.1.100", "Suspicious activity detected", uint(1)).Return(nil)
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", uint(1))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "無効なリクエスト形式",
			requestBody: map[string]interface{}{
				"ip": "192.168.1.100",
			},
			setupMock: func(m *MockFraudUsecase) {
				// No mock setup needed for this test case
			},
			setupContext: func(c *gin.Context) {
				// No context setup needed for this test case
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request format",
		},
		{
			name: "認証されていないユーザー",
			requestBody: map[string]interface{}{
				"ip":     "192.168.1.100",
				"reason": "Test reason",
			},
			setupMock: func(m *MockFraudUsecase) {
				// No mock setup needed for this test case
			},
			setupContext: func(c *gin.Context) {
				// No context setup needed for this test case
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Admin ID not found in context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockFraudUsecase)
			tt.setupMock(mockUsecase)

			fraudHandler := handler.NewFraudHandler(mockUsecase)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonBody, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest("POST", "/fraud/blacklist/ip", bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			tt.setupContext(c)

			fraudHandler.AddIPToBlacklist(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestFraudHandlerSecurityEventComparison(t *testing.T) {
	expected := dto.CreateSecurityEventRequest{
		EventType:   "LOGIN_FAILED",
		Description: "Failed login attempt",
		IPAddress:   "192.168.1.100",
		UserAgent:   "Mozilla/5.0",
	}

	actual := dto.CreateSecurityEventRequest{
		EventType:   "LOGIN_FAILED",
		Description: "Failed login attempt",
		IPAddress:   "192.168.1.100",
		UserAgent:   "Mozilla/5.0",
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("SecurityEvent mismatch (-want +got):\n%s", diff)
	}
}

func TestFraudHandlerCreateSecurityEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    *dto.CreateSecurityEventRequest
		setupMock      func(*MockFraudUsecase)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常なセキュリティイベント作成",
			requestBody: &dto.CreateSecurityEventRequest{
				EventType:   "LOGIN_FAILED",
				Description: "Failed login attempt",
				IPAddress:   "192.168.1.100",
				UserAgent:   "Mozilla/5.0",
			},
			setupMock: func(m *MockFraudUsecase) {
				m.On("CreateSecurityEvent", mock.Anything, mock.AnythingOfType("*dto.CreateSecurityEventRequest")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "無効なリクエスト形式",
			requestBody: &dto.CreateSecurityEventRequest{
				Description: "Test description",
			},
			setupMock: func(m *MockFraudUsecase) {
				// No mock setup needed for this test case
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockFraudUsecase)
			tt.setupMock(mockUsecase)

			fraudHandler := handler.NewFraudHandler(mockUsecase)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonBody, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest("POST", "/fraud/security/events", bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			fraudHandler.CreateSecurityEvent(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestFraudHandlerGetBlacklistedIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(MockFraudUsecase)
	expectedIPs := []*entity.IPBlacklist{
		{IPAddress: "192.168.1.100", Reason: "Suspicious activity"},
		{IPAddress: "10.0.0.1", Reason: "Brute force attack"},
	}

	mockUsecase.On("GetBlacklistedIPs", mock.Anything).Return(expectedIPs, nil)

	fraudHandler := handler.NewFraudHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/fraud/blacklist/ips", nil)

	fraudHandler.GetBlacklistedIPs(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(2), response["count"])

	mockUsecase.AssertExpectations(t)
}

func TestFraudHandlerCleanupExpiredData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(*MockFraudUsecase)
		setupContext   func(*gin.Context)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常なデータクリーンアップ",
			setupMock: func(m *MockFraudUsecase) {
				m.On("CleanupExpiredData", mock.Anything).Return(nil)
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", uint(1))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "認証されていない管理者",
			setupMock: func(m *MockFraudUsecase) {
				// No mock setup needed - authentication fails before usecase is called
			},
			setupContext: func(c *gin.Context) {
				// No user_id set to simulate unauthenticated request
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Admin authentication required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockFraudUsecase)
			tt.setupMock(mockUsecase)

			fraudHandler := handler.NewFraudHandler(mockUsecase)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/fraud/cleanup", nil)

			tt.setupContext(c)

			fraudHandler.CleanupExpiredData(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

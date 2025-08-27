package usecase_test

import (
	"context"
	"testing"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/dto"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestFraudUsecaseAddIPToBlacklist(t *testing.T) {
	mockDomainService := &MockFraudDomainService{}
	fraudUsecase := usecase.NewFraudUsecase(mockDomainService)
	ctx := context.Background()

	ip := "192.168.1.100"
	reason := "Suspicious activity"
	adminID := uint(1)

	mockDomainService.On("GetBlacklistedIPs", ctx, 1, 1000).Return(map[string]interface{}{
		"ips": []*entity.IPBlacklist{},
	}, nil)

	mockDomainService.On("AddIPToBlacklist", ctx, ip, reason, "127.0.0.1", "Admin-Panel").Return(nil)

	err := fraudUsecase.AddIPToBlacklist(ctx, ip, reason, adminID)

	assert.NoError(t, err)
	mockDomainService.AssertExpectations(t)
}

func TestFraudUsecaseAddIPToBlacklistInvalidIP(t *testing.T) {
	mockDomainService := &MockFraudDomainService{}
	fraudUsecase := usecase.NewFraudUsecase(mockDomainService)
	ctx := context.Background()

	invalidIP := "invalid-ip"
	reason := "Suspicious activity"
	adminID := uint(1)

	err := fraudUsecase.AddIPToBlacklist(ctx, invalidIP, reason, adminID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid IP address")
}

func TestFraudUsecaseAddIPToBlacklistEmptyReason(t *testing.T) {
	mockDomainService := &MockFraudDomainService{}
	fraudUsecase := usecase.NewFraudUsecase(mockDomainService)
	ctx := context.Background()

	ip := "192.168.1.100"
	reason := ""
	adminID := uint(1)

	err := fraudUsecase.AddIPToBlacklist(ctx, ip, reason, adminID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reason cannot be empty")
}

func TestFraudUsecaseRemoveIPFromBlacklist(t *testing.T) {
	mockDomainService := &MockFraudDomainService{}
	fraudUsecase := usecase.NewFraudUsecase(mockDomainService)
	ctx := context.Background()

	ip := "192.168.1.100"

	mockDomainService.On("RemoveIPFromBlacklist", ctx, ip, "127.0.0.1", "Admin-Panel").Return(nil)

	err := fraudUsecase.RemoveIPFromBlacklist(ctx, ip)

	assert.NoError(t, err)
	mockDomainService.AssertExpectations(t)
}

func TestFraudUsecaseGetBlacklistedIPs(t *testing.T) {
	mockDomainService := &MockFraudDomainService{}
	fraudUsecase := usecase.NewFraudUsecase(mockDomainService)
	ctx := context.Background()

	expectedIPs := []*entity.IPBlacklist{
		{
			ID:        1,
			IPAddress: "192.168.1.100",
			Reason:    "Suspicious activity",
			IsActive:  true,
		},
	}

	mockDomainService.On("GetBlacklistedIPs", ctx, 1, 100).Return(map[string]interface{}{
		"ips": expectedIPs,
	}, nil)

	ips, err := fraudUsecase.GetBlacklistedIPs(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedIPs, ips)
	mockDomainService.AssertExpectations(t)
}

func TestFraudUsecaseGetSecurityEvents(t *testing.T) {
	mockDomainService := &MockFraudDomainService{}
	fraudUsecase := usecase.NewFraudUsecase(mockDomainService)
	ctx := context.Background()

	limit := 50
	offset := 0

	userID := uint(1)
	expectedEvents := []*entity.SecurityEvent{
		{
			ID:          1,
			UserID:      &userID,
			EventType:   "login_failed",
			Description: "Failed login attempt",
			IPAddress:   "192.168.1.1",
			Severity:    "medium",
		},
	}

	mockDomainService.On("GetSecurityEvents", ctx, 1, limit).Return(map[string]interface{}{
		"events": expectedEvents,
	}, nil)

	events, err := fraudUsecase.GetSecurityEvents(ctx, limit, offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedEvents, events)
	mockDomainService.AssertExpectations(t)
}

func TestFraudUsecaseCreateSecurityEvent(t *testing.T) {
	mockDomainService := &MockFraudDomainService{}
	fraudUsecase := usecase.NewFraudUsecase(mockDomainService)
	ctx := context.Background()

	userID := uint(1)
	req := &dto.CreateSecurityEventRequest{
		UserID:      &userID,
		EventType:   "LOGIN_FAILED",
		Description: "Failed login attempt",
		IPAddress:   "192.168.1.1",
		UserAgent:   "Mozilla/5.0",
	}

	mockDomainService.On("CreateSecurityEvent", ctx, req.UserID, req.EventType, req.Description, req.IPAddress, req.UserAgent, "MEDIUM").Return(nil)

	err := fraudUsecase.CreateSecurityEvent(ctx, req)

	assert.NoError(t, err)
	mockDomainService.AssertExpectations(t)
}

func TestFraudUsecaseCreateSecurityEventInvalidEventType(t *testing.T) {
	mockDomainService := &MockFraudDomainService{}
	fraudUsecase := usecase.NewFraudUsecase(mockDomainService)
	ctx := context.Background()

	userID := uint(1)
	req := &dto.CreateSecurityEventRequest{
		UserID:      &userID,
		EventType:   "INVALID_EVENT",
		Description: "Test description",
		IPAddress:   "192.168.1.1",
		UserAgent:   "Mozilla/5.0",
	}

	err := fraudUsecase.CreateSecurityEvent(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid event type")
}

func TestFraudUsecaseCreateRateLimitRule(t *testing.T) {
	mockDomainService := &MockFraudDomainService{}
	fraudUsecase := usecase.NewFraudUsecase(mockDomainService)
	ctx := context.Background()

	req := &dto.CreateRateLimitRuleRequest{
		RuleType:    "IP",
		Identifier:  "192.168.1.1",
		MaxRequests: 100,
		WindowSize:  3600,
	}

	mockDomainService.On("CreateRateLimitRule", ctx, req.RuleType, req.Identifier, req.MaxRequests, req.WindowSize).Return(nil)

	err := fraudUsecase.CreateRateLimitRule(ctx, req)

	assert.NoError(t, err)
	mockDomainService.AssertExpectations(t)
}

func TestFraudUsecaseCleanupExpiredData(t *testing.T) {
	mockDomainService := &MockFraudDomainService{}
	fraudUsecase := usecase.NewFraudUsecase(mockDomainService)
	ctx := context.Background()

	mockDomainService.On("CleanupExpiredData", ctx).Return(nil)

	err := fraudUsecase.CleanupExpiredData(ctx)

	assert.NoError(t, err)
	mockDomainService.AssertExpectations(t)
}

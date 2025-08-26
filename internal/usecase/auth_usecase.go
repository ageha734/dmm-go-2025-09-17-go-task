package usecase

import (
	"context"
	"fmt"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/service"
)

type AuthUsecase struct {
	authDomainService  service.AuthDomainServiceInterface
	fraudDomainService service.FraudDomainServiceInterface
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	User         *entity.User `json:"user"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Age      int    `json:"age" binding:"min=0,max=150"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

func NewAuthUsecase(authDomainService service.AuthDomainServiceInterface, fraudDomainService service.FraudDomainServiceInterface) *AuthUsecase {
	return &AuthUsecase{
		authDomainService:  authDomainService,
		fraudDomainService: fraudDomainService,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, req RegisterRequest, ipAddress, userAgent string) (*LoginResponse, error) {
	fraudAnalysis, err := u.fraudDomainService.AnalyzeFraud(ctx, nil, req.Email, ipAddress, userAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze fraud: %w", err)
	}

	if fraudAnalysis.IsHighRisk() {
		_ = u.fraudDomainService.CreateSecurityEvent(ctx, nil, "HIGH_RISK_REGISTRATION",
			"High risk registration attempt", ipAddress, userAgent, "HIGH")
		return nil, fmt.Errorf("registration blocked due to security concerns")
	}

	user, err := u.authDomainService.Register(ctx, req.Name, req.Email, req.Password, req.Age)
	if err != nil {
		_ = u.fraudDomainService.RecordLoginAttempt(ctx, req.Email, ipAddress, userAgent, false, "Registration failed")
		return nil, err
	}

	_ = u.fraudDomainService.RecordLoginAttempt(ctx, req.Email, ipAddress, userAgent, true, "")

	_ = u.fraudDomainService.CreateSecurityEvent(ctx, &user.ID, "USER_REGISTRATION",
		"New user registered", ipAddress, userAgent, "LOW")

	auth, roles, err := u.authDomainService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return &LoginResponse{User: user}, nil
	}

	accessToken, err := u.authDomainService.GenerateAccessToken(auth.UserID, auth.Email, roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := u.authDomainService.GenerateRefreshToken(ctx, auth.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
		User:         user,
	}, nil
}

func (u *AuthUsecase) Login(ctx context.Context, req LoginRequest, ipAddress, userAgent string) (*LoginResponse, error) {
	fraudAnalysis, err := u.fraudDomainService.AnalyzeFraud(ctx, nil, req.Email, ipAddress, userAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze fraud: %w", err)
	}

	if fraudAnalysis.IsHighRisk() {
		_ = u.fraudDomainService.RecordLoginAttempt(ctx, req.Email, ipAddress, userAgent, false, "High risk login blocked")
		_ = u.fraudDomainService.CreateSecurityEvent(ctx, nil, "HIGH_RISK_LOGIN",
			"High risk login attempt blocked", ipAddress, userAgent, "HIGH")
		return nil, fmt.Errorf("login blocked due to security concerns")
	}

	auth, roles, err := u.authDomainService.Login(ctx, req.Email, req.Password)
	if err != nil {
		_ = u.fraudDomainService.RecordLoginAttempt(ctx, req.Email, ipAddress, userAgent, false, err.Error())
		return nil, err
	}

	_ = u.fraudDomainService.RecordLoginAttempt(ctx, req.Email, ipAddress, userAgent, true, "")

	_ = u.fraudDomainService.CreateSecurityEvent(ctx, &auth.UserID, "LOGIN",
		"User logged in successfully", ipAddress, userAgent, "LOW")

	accessToken, err := u.authDomainService.GenerateAccessToken(auth.UserID, auth.Email, roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := u.authDomainService.GenerateRefreshToken(ctx, auth.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	user := &entity.User{
		ID:    auth.UserID,
		Email: auth.Email,
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
		User:         user,
	}, nil
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*LoginResponse, error) {
	auth, roles, newRefreshToken, err := u.authDomainService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.authDomainService.GenerateAccessToken(auth.UserID, auth.Email, roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	user := &entity.User{
		ID:    auth.UserID,
		Email: auth.Email,
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600,
		User:         user,
	}, nil
}

func (u *AuthUsecase) ChangePassword(ctx context.Context, userID uint, req ChangePasswordRequest, ipAddress, userAgent string) error {
	err := u.authDomainService.ChangePassword(ctx, userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return err
	}

	_ = u.fraudDomainService.CreateSecurityEvent(ctx, &userID, "PASSWORD_CHANGE",
		"User changed password", ipAddress, userAgent, "MEDIUM")

	return nil
}

func (u *AuthUsecase) Logout(ctx context.Context, userID uint, sessionID, ipAddress, userAgent string) error {
	if err := u.authDomainService.Logout(ctx, userID); err != nil {
		return err
	}

	if sessionID != "" {
		_ = u.fraudDomainService.DeactivateUserSession(ctx, sessionID)
	}

	_ = u.fraudDomainService.CreateSecurityEvent(ctx, &userID, "LOGOUT",
		"User logged out", ipAddress, userAgent, "LOW")

	return nil
}

func (u *AuthUsecase) ValidateToken(tokenString string) (*service.JWTClaims, error) {
	return u.authDomainService.ValidateToken(tokenString)
}

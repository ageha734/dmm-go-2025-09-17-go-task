package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
)

type AuthDomainService struct {
	userRepo         repository.UserRepository
	authRepo         repository.AuthRepository
	roleRepo         repository.RoleRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtSecret        string
}

type JWTClaims struct {
	UserID uint     `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

func NewAuthDomainService(
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
	roleRepo repository.RoleRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtSecret string,
) *AuthDomainService {
	return &AuthDomainService{
		userRepo:         userRepo,
		authRepo:         authRepo,
		roleRepo:         roleRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        jwtSecret,
	}
}

func (s *AuthDomainService) Register(ctx context.Context, name, email, password string, age int) (*entity.User, error) {
	tempAuth := &entity.Auth{}
	if err := tempAuth.ValidatePassword(password); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	exists, err := s.authRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	user := entity.NewUser(name, email, age)
	if !user.IsValidAge() || !user.IsValidEmail() {
		return nil, errors.New("invalid user data")
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	auth, err := entity.NewAuth(user.ID, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth: %w", err)
	}

	if err := s.authRepo.Create(ctx, auth); err != nil {
		return nil, fmt.Errorf("failed to create auth: %w", err)
	}

	if err := s.assignDefaultRole(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to assign default role: %w", err)
	}

	return user, nil
}

func (s *AuthDomainService) Login(ctx context.Context, email, password string) (*entity.Auth, []string, error) {
	auth, err := s.authRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if !auth.IsActive {
		return nil, nil, ErrInvalidCredentials
	}

	if err := auth.VerifyPassword(password); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	auth.UpdateLastLogin()
	if err := s.authRepo.Update(ctx, auth); err != nil {
		return nil, nil, fmt.Errorf("failed to update auth: %w", err)
	}

	roles, err := s.roleRepo.GetUserRoles(ctx, auth.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	return auth, roleNames, nil
}

func (s *AuthDomainService) GenerateAccessToken(userID uint, email string, roles []string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dmm-go-task",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthDomainService) GenerateRefreshToken(ctx context.Context, userID uint) (string, error) {
	tokenStr := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	refreshToken := entity.NewRefreshToken(userID, tokenStr, expiresAt)
	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return tokenStr, nil
}

func (s *AuthDomainService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (s *AuthDomainService) RefreshToken(ctx context.Context, refreshTokenStr string) (*entity.Auth, []string, string, error) {
	refreshToken, err := s.refreshTokenRepo.GetByToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, nil, "", ErrInvalidToken
	}

	if !refreshToken.IsValid() {
		return nil, nil, "", ErrInvalidToken
	}

	auth, err := s.authRepo.GetByUserID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, nil, "", ErrUserNotFound
	}

	roles, err := s.roleRepo.GetUserRoles(ctx, auth.UserID)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	refreshToken.Revoke()
	if err := s.refreshTokenRepo.Update(ctx, refreshToken); err != nil {
		return nil, nil, "", fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	newRefreshToken, err := s.GenerateRefreshToken(ctx, auth.UserID)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	return auth, roleNames, newRefreshToken, nil
}

func (s *AuthDomainService) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	auth, err := s.authRepo.GetByUserID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := auth.ChangePassword(currentPassword, newPassword); err != nil {
		return err
	}

	if err := s.authRepo.Update(ctx, auth); err != nil {
		return fmt.Errorf("failed to update auth: %w", err)
	}

	if err := s.refreshTokenRepo.RevokeByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke refresh tokens: %w", err)
	}

	return nil
}

func (s *AuthDomainService) Logout(ctx context.Context, userID uint) error {
	return s.refreshTokenRepo.RevokeByUserID(ctx, userID)
}

func (s *AuthDomainService) assignDefaultRole(ctx context.Context, userID uint) error {
	role, err := s.roleRepo.GetByName(ctx, "user")
	if err != nil {
		defaultRole := entity.NewRole("user", "Default user role")
		if err := s.roleRepo.Create(ctx, defaultRole); err != nil {
			return fmt.Errorf("failed to create default role: %w", err)
		}
		role = defaultRole
	}

	return s.roleRepo.AssignToUser(ctx, userID, role.ID)
}

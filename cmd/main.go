package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/handler"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/middleware"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/service"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/infrastructure/external"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/infrastructure/persistence"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := initDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	jwtSecret := getJWTSecret()
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	userRepo := persistence.NewUserRepository(db)
	authRepo := persistence.NewAuthRepository(db)
	roleRepo := persistence.NewRoleRepository(db)
	refreshTokenRepo := persistence.NewRefreshTokenRepository(db)
	userProfileRepo := persistence.NewUserProfileRepository(db)
	userMembershipRepo := persistence.NewUserMembershipRepository(db)
	securityEventRepo := persistence.NewSecurityEventRepository(db)
	ipBlacklistRepo := persistence.NewIPBlacklistRepository(db)
	loginAttemptRepo := persistence.NewLoginAttemptRepository(db)
	rateLimitRuleRepo := persistence.NewRateLimitRuleRepository(db)
	userSessionRepo := persistence.NewUserSessionRepository(db)
	deviceFingerprintRepo := persistence.NewDeviceFingerprintRepository(db)

	authDomainService := service.NewAuthDomainService(
		userRepo,
		authRepo,
		roleRepo,
		refreshTokenRepo,
		jwtSecret,
	)

	fraudDomainService := service.NewFraudDomainService(
		securityEventRepo,
		ipBlacklistRepo,
		loginAttemptRepo,
		rateLimitRuleRepo,
		userSessionRepo,
		deviceFingerprintRepo,
	)

	authUsecase := usecase.NewAuthUsecase(authDomainService, fraudDomainService)
	userUsecase := usecase.NewUserUsecase(
		userRepo,
		userProfileRepo,
		userMembershipRepo,
		fraudDomainService,
	)
	fraudUsecase := usecase.NewFraudUsecase(fraudDomainService)

	redisClient := external.NewRedisClient("localhost:6379", "", 0)
	cacheService := external.NewCacheService(redisClient)

	authMiddleware := middleware.NewAuthMiddleware(authDomainService, cacheService)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(cacheService)

	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	fraudHandler := handler.NewFraudHandler(fraudUsecase)

	router := setupRouter(authHandler, userHandler, fraudHandler, authMiddleware, rateLimitMiddleware)

	port := getPort()
	log.Printf("Starting server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDatabase() (*gorm.DB, error) {
	dsn := getDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupRouter(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, fraudHandler *handler.FraudHandler, authMiddleware *middleware.AuthMiddleware, rateLimitMiddleware *middleware.RateLimitMiddleware) *gin.Engine {
	router := gin.Default()

	router.Use(handler.CORSMiddleware())
	router.Use(handler.SecurityHeadersMiddleware())
	router.Use(handler.RequestIDMiddleware())

	router.HandleMethodNotAllowed = true
	router.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{"error": "Method not allowed"})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Not found"})
	})

	router.GET("/health", userHandler.HealthCheck)

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.Use(rateLimitMiddleware.RateLimitByIP(10, time.Minute))
			auth.POST("/register", rateLimitMiddleware.RateLimitByIP(5, time.Minute), authHandler.Register)
			auth.POST("/login", rateLimitMiddleware.RateLimitByIP(5, time.Minute), authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/validate", authHandler.ValidateToken)
		}

		user := v1.Group("/user")
		user.Use(authMiddleware.RequireAuth())
		{
			user.POST("/logout", authHandler.Logout)
			user.POST("/change-password", authHandler.ChangePassword)
			user.GET("/profile", userHandler.GetUserProfile)
			user.PUT("/profile", authHandler.UpdateUserProfile)
			user.GET("/dashboard", authHandler.GetUserDashboard)
			user.GET("/notifications", authHandler.GetUserNotifications)
			user.PUT("/notifications/:id/read", authHandler.MarkNotificationRead)
			user.GET("/points/transactions", authHandler.GetUserPointTransactions)
			user.POST("/preferences", authHandler.SetUserPreference)
			user.GET("/preferences", authHandler.GetUserPreferences)
		}

		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)

			users.PUT("/:id", authMiddleware.RequireAuth(), userHandler.UpdateUser)
			users.DELETE("/:id", authMiddleware.RequireAuth(), userHandler.DeleteUser)
		}

		admin := v1.Group("/admin")
		admin.Use(authMiddleware.RequireAuth(), authMiddleware.RequireRole("admin"))
		{
			admin.GET("/health", userHandler.GetSystemHealth)
			admin.GET("/users", userHandler.GetUsers)
			admin.GET("/users/:user_id", userHandler.GetUserDetails)
			admin.POST("/users/:user_id/points", userHandler.AddPointsToUser)
			admin.POST("/users/:user_id/notifications", userHandler.CreateNotificationForUser)
			admin.POST("/points/expire", userHandler.ExpireUserPoints)
		}

		fraud := v1.Group("/fraud")
		fraud.Use(authMiddleware.RequireAuth(), authMiddleware.RequireRole("admin"))
		{
			fraud.POST("/blacklist/ip", fraudHandler.AddIPToBlacklist)
			fraud.DELETE("/blacklist/ip/:ip", fraudHandler.RemoveIPFromBlacklist)
			fraud.GET("/blacklist/ips", fraudHandler.GetBlacklistedIPs)

			fraud.GET("/security/events", fraudHandler.GetSecurityEvents)
			fraud.POST("/security/events", fraudHandler.CreateSecurityEvent)
			fraud.POST("/security/blacklist", fraudHandler.AddIPToBlacklist)
			fraud.DELETE("/security/blacklist/:ip", fraudHandler.RemoveIPFromBlacklist)

			fraud.POST("/rate/limits", fraudHandler.CreateRateLimitRule)
			fraud.PUT("/rate/limits/:id", fraudHandler.UpdateRateLimitRule)
			fraud.DELETE("/rate/limits/:id", fraudHandler.DeleteRateLimitRule)
			fraud.GET("/rate/limits", fraudHandler.GetRateLimitRules)

			fraud.GET("/sessions", fraudHandler.GetActiveSessions)
			fraud.DELETE("/sessions/:sessionId", fraudHandler.DeactivateSession)

			fraud.GET("/devices", fraudHandler.GetDevices)
			fraud.PUT("/devices/:fingerprint/trust", fraudHandler.TrustDevice)

			fraud.POST("/cleanup", fraudHandler.CleanupExpiredData)
		}

		v1.GET("/stats", userHandler.GetUserStats)
	}

	return router
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]
		jwtSecret := getJWTSecret()

		claims, err := validateJWTToken(token, jwtSecret)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		c.Next()
	}
}

func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("user_roles")
		if !exists {
			c.JSON(403, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		roleSlice, ok := roles.([]string)
		if !ok {
			c.JSON(403, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		hasAdminRole := false
		for _, role := range roleSlice {
			if role == "admin" {
				hasAdminRole = true
				break
			}
		}

		if !hasAdminRole {
			c.JSON(403, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func validateJWTToken(tokenString, jwtSecret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

type JWTClaims struct {
	UserID uint     `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

func getDSN() string {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "testuser:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return dsn
}

func getJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

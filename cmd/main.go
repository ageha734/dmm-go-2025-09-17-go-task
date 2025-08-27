package main

import (
	"log"
	"os"
	"strconv"
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

	redisClient := external.NewRedisClient(getRedisAddr(), getRedisPassword(), getRedisDB())
	cacheService := external.NewCacheService(redisClient)

	authDomainService := service.NewAuthDomainService(
		userRepo,
		authRepo,
		roleRepo,
		refreshTokenRepo,
		cacheService,
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

	authUsecase := usecase.NewAuthUsecase(authDomainService, fraudDomainService, cacheService)
	userUsecase := usecase.NewUserUsecase(
		userRepo,
		userProfileRepo,
		userMembershipRepo,
		fraudDomainService,
		redisClient,
	)
	fraudUsecase := usecase.NewFraudUsecase(fraudDomainService)

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
	router.Use(middleware.ContentTypeMiddleware())

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
			auth.Use(rateLimitMiddleware.RateLimitByIP(getAuthRateLimit(), time.Minute))
			auth.POST("/register", rateLimitMiddleware.RateLimitByIP(getRegisterRateLimit(), time.Minute), authHandler.Register)
			auth.POST("/login", rateLimitMiddleware.RateLimitByIP(getLoginRateLimit(), time.Minute), authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/validate", authHandler.ValidateToken)
		}

		user := v1.Group("/user")
		user.Use(authMiddleware.RequireAuth())
		{
			user.POST("/logout", authHandler.Logout)
			user.POST("/change-password", authHandler.ChangePassword)
			user.GET("/profile", userHandler.GetUserProfile)
			user.PUT("/profile/:id", userHandler.UpdateUserProfile)
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

type JWTClaims struct {
	UserID uint     `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

func getDSN() string {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "testuser:password@tcp(mysql:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
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

func getRedisAddr() string {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "redis"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	return host + ":" + port
}

func getRedisPassword() string {
	password := os.Getenv("REDIS_PASSWORD")
	if password == "" {
		password = "password"
	}
	return password
}

func getRedisDB() int {
	return 0
}

func getAuthRateLimit() int64 {
	limit := os.Getenv("AUTH_RATE_LIMIT")
	if limit == "" {
		return 100
	}
	val, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return 100
	}
	return val
}

func getRegisterRateLimit() int64 {
	limit := os.Getenv("REGISTER_RATE_LIMIT")
	if limit == "" {
		return 50
	}
	val, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return 50
	}
	return val
}

func getLoginRateLimit() int64 {
	limit := os.Getenv("LOGIN_RATE_LIMIT")
	if limit == "" {
		return 50
	}
	val, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return 50
	}
	return val
}

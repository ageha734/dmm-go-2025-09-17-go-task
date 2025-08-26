package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/handler"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/service"
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

	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	router := setupRouter(authHandler, userHandler)

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

func setupRouter(authHandler *handler.AuthHandler, userHandler *handler.UserHandler) *gin.Engine {
	router := gin.Default()

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
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/validate", authHandler.ValidateToken)
		}

		authenticated := v1.Group("")
		authenticated.Use(authMiddleware())
		{
			authenticated.POST("/auth/logout", authHandler.Logout)
			authenticated.POST("/auth/change-password", authHandler.ChangePassword)
		}

		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)

			users.PUT("/:id", authMiddleware(), userHandler.UpdateUser)
			users.DELETE("/:id", authMiddleware(), userHandler.DeleteUser)
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
		dsn = "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
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

package routes

import (
	"github.com/dmm-com/dmm-go-2025-09-17-go-task/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", handlers.HealthCheck)

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("", handlers.GetUsers)          // GET /api/v1/users
			users.POST("", handlers.CreateUser)       // POST /api/v1/users
			users.GET("/:id", handlers.GetUser)       // GET /api/v1/users/:id
			users.PUT("/:id", handlers.UpdateUser)    // PUT /api/v1/users/:id
			users.DELETE("/:id", handlers.DeleteUser) // DELETE /api/v1/users/:id
		}

		v1.GET("/stats", handlers.GetUserStats) // GET /api/v1/stats
	}

	return router
}

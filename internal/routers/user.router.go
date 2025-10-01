package routers

import (
	"github.com/Darari17/social-media/internal/handlers"
	"github.com/Darari17/social-media/internal/middlewares"
	"github.com/Darari17/social-media/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitUserRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	user := router.Group("/users")
	userRepo := repos.NewUserRepo(db)
	userHandler := handlers.NewUserHandler(userRepo)

	user.GET("", userHandler.GetAllUsers)
	user.GET("/profile", middlewares.RequiredToken(rdb), userHandler.GetUserByID)
	user.PATCH("/profile", middlewares.RequiredToken(rdb), userHandler.UpdateUser)
}

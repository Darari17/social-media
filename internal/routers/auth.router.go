package routers

import (
	"github.com/Darari17/social-media/internal/handlers"
	"github.com/Darari17/social-media/internal/middlewares"
	"github.com/Darari17/social-media/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitAuthRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	auth := r.Group("/auth")
	authRepo := repos.NewAuthRepo(db, rdb)
	authHandler := handlers.NewAuthHandler(authRepo)

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.DELETE("/logout", middlewares.RequiredToken(rdb), authHandler.Logout)
}

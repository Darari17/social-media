package routers

import (
	"github.com/Darari17/social-media/internal/handlers"
	"github.com/Darari17/social-media/internal/middlewares"
	"github.com/Darari17/social-media/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitFollowRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	followRepo := repos.NewFollowRepo(db)
	followHandler := handlers.NewFollowHandler(followRepo)

	follow := router.Group("/follow")
	follow.POST("/:id", middlewares.RequiredToken(rdb), followHandler.FollowUser)
	follow.DELETE("/:id", middlewares.RequiredToken(rdb), followHandler.UnfollowUser)

	users := router.Group("/users")
	users.GET("/:id/followers", followHandler.GetFollowers)
	users.GET("/:id/following", followHandler.GetFollowing)
}

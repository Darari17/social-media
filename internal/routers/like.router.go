package routers

import (
	"github.com/Darari17/social-media/internal/handlers"
	"github.com/Darari17/social-media/internal/middlewares"
	"github.com/Darari17/social-media/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitLikeRoutes(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	likeRepo := repos.NewLikeRepo(db)
	postRepo := repos.NewPostRepo(db)

	likeHandler := handlers.NewLikeHandler(likeRepo, postRepo)

	posts := r.Group("/posts")
	posts.POST("/:id/like", middlewares.RequiredToken(rdb), likeHandler.LikePost)
	posts.DELETE("/:id/like", middlewares.RequiredToken(rdb), likeHandler.UnlikePost)
	posts.GET("/:id/likes", middlewares.RequiredToken(rdb), likeHandler.GetLikes)
}

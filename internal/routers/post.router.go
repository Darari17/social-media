package routers

import (
	"github.com/Darari17/social-media/internal/handlers"
	"github.com/Darari17/social-media/internal/middlewares"
	"github.com/Darari17/social-media/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitPostRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	postRepo := repos.NewPostRepo(db, rdb)
	postHandler := handlers.NewPostHandler(postRepo)

	posts := router.Group("/posts")

	posts.POST("", middlewares.RequiredToken(rdb), postHandler.CreatePost)
	posts.GET("", postHandler.GetAllPosts)
	posts.GET("/:id", postHandler.GetPostByID)
	posts.PATCH("/:id", middlewares.RequiredToken(rdb), postHandler.UpdatePost)
	posts.DELETE("/:id", middlewares.RequiredToken(rdb), postHandler.DeletePost)
}

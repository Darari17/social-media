package routers

import (
	"github.com/Darari17/social-media/internal/handlers"
	"github.com/Darari17/social-media/internal/middlewares"
	"github.com/Darari17/social-media/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitCommentRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	commentRepo := repos.NewCommentRepo(db)
	postRepo := repos.NewPostRepo(db, rdb)
	commentHandler := handlers.NewCommentHandler(commentRepo, postRepo)

	post := r.Group("/posts")
	post.POST("/:id/comments", middlewares.RequiredToken(rdb), commentHandler.CreateComment)
	post.GET("/:id/comments", middlewares.RequiredToken(rdb), commentHandler.GetComments)
	post.PUT("/comments/:id", middlewares.RequiredToken(rdb), commentHandler.UpdateComment)
	post.DELETE("/comments/:id", middlewares.RequiredToken(rdb), commentHandler.DeleteComment)
}

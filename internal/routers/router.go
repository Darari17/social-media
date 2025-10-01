package routers

import (
	"net/http"

	"github.com/Darari17/social-media/internal/dtos"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitRouter(db *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	InitAuthRouter(r, db, rdb)
	InitUserRouter(r, db, rdb)
	InitPostRouter(r, db, rdb)
	InitFollowRouter(r, db, rdb)
	InitLikeRoutes(r, db, rdb)

	r.Static("/img", "public")

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Page not found",
		})
	})

	return r
}

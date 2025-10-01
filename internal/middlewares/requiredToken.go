package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/Darari17/social-media/internal/dtos"
	"github.com/Darari17/social-media/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func RequiredToken(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		if bearerToken == "" {
			log.Println("Authorization header is missing")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Response{
				Code:    http.StatusUnauthorized,
				Success: false,
				Message: "Please log in first",
			})
			return
		}

		parts := strings.Split(bearerToken, " ")
		if len(parts) != 2 {
			log.Println("Invalid Authorization header format. Expected: Bearer <token>")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Response{
				Code:    http.StatusUnauthorized,
				Success: false,
				Message: "Format authorization header invalid",
			})
			return
		}

		if parts[0] != "Bearer" {
			log.Println("Authorization header must start with 'Bearer'")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Response{
				Code:    http.StatusUnauthorized,
				Success: false,
				Message: "Format authorization header invalid",
			})
			return
		}

		token := parts[1]
		if token == "" {
			log.Println("Token is empty after Bearer")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Response{
				Code:    http.StatusUnauthorized,
				Success: false,
				Message: "Please log in first",
			})
			return
		}

		isBlacklist, err := rdb.Get(ctx, "Mosting:blacklist:"+token).Result()
		if err == nil && isBlacklist == "true" {
			log.Println("The token has logged out, please log in again")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Response{
				Code:    http.StatusUnauthorized,
				Success: false,
				Message: "The token has logged out, please log in again",
			})
			return
		} else if err != redis.Nil && err != nil {
			log.Println("Error when checking blacklist redis cache:", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dtos.Response{
				Code:    http.StatusInternalServerError,
				Success: false,
				Message: "Internal server error",
			})
			return
		}

		claims := &pkg.Claims{}
		if err := claims.VerifyToken(token); err != nil {
			if strings.Contains(err.Error(), jwt.ErrTokenInvalidIssuer.Error()) {
				log.Println("JWT Error.\nCause: ", err.Error())
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Response{
					Code:    http.StatusUnauthorized,
					Success: false,
					Message: "Please log in again",
				})
				return
			}
			if strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
				log.Println("JWT Error.\nCause: ", err.Error())
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Response{
					Code:    http.StatusUnauthorized,
					Success: false,
					Message: "Please log in again",
				})
				return
			}

			log.Println(jwt.ErrTokenExpired)
			log.Println("Internal Server Error.\nCause: ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dtos.Response{
				Code:    http.StatusInternalServerError,
				Success: false,
				Message: "Internal server error",
			})
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

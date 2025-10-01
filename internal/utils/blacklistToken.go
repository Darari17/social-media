package utils

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func BlackListTokenRedish(c context.Context, rdb redis.Client, token string) error {
	err := rdb.Set(c, "Mosting:blacklist:"+token, "true", 30*time.Minute).Err()
	if err != nil {
		log.Println("Redis Error when blacklist token:", err)
		return err
	}
	return nil
}

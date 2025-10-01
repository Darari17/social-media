package main

import (
	"log"

	"github.com/Darari17/social-media/internal/configs"
	"github.com/Darari17/social-media/internal/routers"
	"github.com/joho/godotenv"
)

// @title 											Social Media
// @version 										1.0
// @securityDefinitions.apikey 	BearerAuth
// @in 													header
// @name 												Authorization
// @description									RESTful API created using gin for Backend Social media
// @host												localhost:8080
// @basePath										/
func main() {
	// load env
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load env.\nCause:", err.Error())
		return
	}

	// init db
	db, err := configs.InitDB()
	if err != nil {
		log.Println("Failed to connect to database.\nCause:", err.Error())
		return
	}
	log.Println("DB Connected")
	defer db.Close()

	// test db
	if err := configs.PingDB(db); err != nil {
		log.Println("Ping to DB failed.\nCause:", err.Error())
		return
	}

	// init redis
	rdb, err := configs.InitRedis()
	if err != nil {
		log.Println("Failed to connect Redis.\nCause: ", err.Error())
		return
	}
	log.Println("Redis Connected.")
	defer rdb.Close()

	// router
	router := routers.InitRouter(db, rdb)
	router.Run(":8080")
}

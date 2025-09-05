package server

import (
	"Badminton-Hub/internal/adapter/inbound/handler/http/gin"
	redisCache "Badminton-Hub/internal/adapter/outbound/cache/redis"
	mongodb "Badminton-Hub/internal/adapter/outbound/database/mongoDB"
	"Badminton-Hub/internal/core/service"
	"Badminton-Hub/util"
	"log"
)

func StartServer() {
	defer util.ShutdownServer()

	// Load configuration
	err := util.SetConfig()
	if err != nil {
		log.Fatalln("Failed to load configuration: " + err.Error())
	}
	config := util.LoadConfig()

	// Initialize MongoDB
	db := mongodb.NewMongoDB(config.DBName)

	// Initialize Redis cache
	cacheRedis := redisCache.NewRedisCache()

	// Initialize services
	encryptionJWT := service.NewJWTEncryption()
	middleware := service.NewMiddlewareUtil(encryptionJWT, cacheRedis)
	memberUtil := service.NewMemberUtil(db, middleware)
	redirectUtil := service.NewRedirectUtil(cacheRedis)

	// Initialize HTTP server
	externalRoute := gin.NewGinMainRoute(
		middleware,
		memberUtil,
		redirectUtil,
	)

	externalRoute.Start()
	defer externalRoute.Run()

	externalRoute.RouteMember()
	externalRoute.RouteTest()
}

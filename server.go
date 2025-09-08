package server

import (
	"Badminton-Hub/internal/adapter/inbound/handler/http/gin"
	redisCache "Badminton-Hub/internal/adapter/outbound/cache/redis"
	mongodb "Badminton-Hub/internal/adapter/outbound/database/mongoDB"
	"Badminton-Hub/internal/core/service"
	core_util "Badminton-Hub/internal/util"
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

	// Initialize Util
	encryptionJWT := core_util.NewJWTEncryptionUtil()

	// Initialize services
	middleware := service.NewMiddlewareService(encryptionJWT, cacheRedis)
	memberUtil := service.NewMemberService(db, middleware)
	redirectUtil := service.NewRedirectService(cacheRedis)

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

package server

import (
	"Badminton-Hub/internal/adapter/inbound/handler/http/gin"
	"Badminton-Hub/internal/adapter/outbound/3rdParty/google"
	"Badminton-Hub/internal/adapter/outbound/cache/redis"
	"Badminton-Hub/internal/adapter/outbound/database/mongoDB"
	"Badminton-Hub/internal/core/service"
	"Badminton-Hub/internal/core_util"
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
	db := mongoDB.NewMongoDB(config.DBName)

	// Initialize Redis cache
	cacheRedis := redis.NewRedisCache()

	// Setup Util
	encryptionJWT := core_util.NewJWTEncryptionUtil()

	// Setup Google
	googleRedirect := google.NewGoogleRedirect(cacheRedis)
	googleCallback := google.NewGoogleCallback(cacheRedis)

	// Initialize services
	callback := service.NewCallbackService(googleCallback)
	middleware := service.NewMiddlewareService(db, encryptionJWT, callback)
	authenticate := service.NewAuthenticationService(db, middleware)

	authenticationSystem := service.NewAuthenticationSystem(authenticate, middleware)
	redirect := service.NewRedirect(googleRedirect)

	// Initialize HTTP server
	externalRoute := gin.NewGinRoute(
		authenticationSystem,
		redirect,
	)

	// Initialize HTTP server
	externalRoute.Start()
	defer externalRoute.Run()

	externalRoute.RouteAuthenticationSystem()
	externalRoute.RouteRedirect()
	externalRoute.RouteCallback()
}

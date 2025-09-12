package server

import (
	"Badminton-Hub/internal/adapter/inbound/handler/http/gin"
	third_party "Badminton-Hub/internal/adapter/outbound/3rdParty"
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
	middlewareUtil := service.NewMiddlewareUtil(encryptionJWT)

	// Setup 3rd Party
	thirdPartyUtil := third_party.NewThirdPartyUtil()
	authenticate3rdParty := third_party.New3rdPartyMiddleware(cacheRedis)
	redirect3rdParty := third_party.New3rdPartyRedirect(cacheRedis)

	// Initialize services
	authenticateUtil := service.NewAuthenticateService(authenticate3rdParty, middlewareUtil, db)
	authentication := service.NewAuthenticationService(db, middlewareUtil, thirdPartyUtil)
	middlewareSystem := service.NewMiddlewareSystem(authenticateUtil, middlewareUtil)
	authenticationSystem := service.NewAuthenticationSystem(authentication, middlewareSystem)
	redirect := service.NewRedirect(redirect3rdParty)
	member := service.NewMemberService(db)

	// Initialize HTTP server
	externalRoute := gin.NewGinRoute(
		authenticationSystem,
		redirect,
		member,
	)

	// Initialize HTTP server
	externalRoute.Start()
	defer externalRoute.Run()

	externalRoute.RouteAuthenticationSystem()
	externalRoute.RouteRedirect()
	externalRoute.RouteCallback()
	externalRoute.RouteMember()
}

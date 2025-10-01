package server

import (
	"Badminton-Hub/internal/adapter/inbound/handler/http/gin"
	third_party "Badminton-Hub/internal/adapter/outbound/3rdParty"
	"Badminton-Hub/internal/adapter/outbound/cache/redis"
	"Badminton-Hub/internal/adapter/outbound/database/mongoDB"
	"Badminton-Hub/internal/adapter/outbound/observability/log/zeroLog"
	"Badminton-Hub/internal/adapter/outbound/observability/metrics/prometheus"
	"Badminton-Hub/internal/core/service"
	"Badminton-Hub/internal/core_util"
	"Badminton-Hub/util"
	"log"
)

func StartServer() {
	// Load configuration
	err := util.SetConfig()
	if err != nil {
		log.Fatalln("Failed to load configuration: " + err.Error())
	}

	// Initialize MongoDB
	db, closeDB := mongoDB.NewMongoDB()

	// Initialize Redis cache
	cacheRedis, closeRedis := redis.NewRedisCache()

	// Graceful shutdown
	defer util.ShutdownServer(closeDB, closeRedis)

	// Setup Observability
	metrics := prometheus.NewPrometheus()
	log := zeroLog.NewZeroLog()

	// Setup Util
	observabilityUtil := core_util.NewObservability(metrics, log)
	encryptionJWT := core_util.NewJWTEncryptionUtil()
	middlewareUtil := core_util.NewMiddlewareUtil(encryptionJWT)

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
	route := gin.NewGinRoute(
		authenticationSystem,
		redirect,
		member,
		observabilityUtil,
	)

	// Initialize HTTP server
	runServer := route.Start()
	defer runServer()

	route.RouteHealthCheck()
	route.RouteAuthenticationSystem()
	route.RouteRedirect()
	route.RouteCallback()
	route.RouteMember()
	route.RouteObservability()
}

package server

import (
	"Badminton-Hub/internal/adapter/inbound/handler/http/gin"
	mongodb "Badminton-Hub/internal/adapter/outbound/database/mongoDB"
	"Badminton-Hub/internal/core/service"
	"Badminton-Hub/util"
)

func StartServer() {
	config, err := util.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize MongoDB
	db := mongodb.NewMongoDB(config.DBName)

	encryptionJWT := service.NewJWTEncryption()
	middleware := service.NewMiddlewareUtil(encryptionJWT)
	memberUtil := service.NewMemberUtil(db, middleware)

	externalRoute := gin.NewGinMainRoute(middleware, memberUtil)

	mainRoute := service.NewMainRoute(externalRoute)

	mainRoute.RouteMember()
}

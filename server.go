package server

import (
	"Badminton-Hub/internal/adapter/inbound/handler/http/gin"
	mongodb "Badminton-Hub/internal/adapter/outbound/database/mongoDB"
	"Badminton-Hub/internal/core/service"
	core_util "Badminton-Hub/internal/util"
	"os"

	"github.com/joho/godotenv"
)

func StartServer() {
	godotenv.Load()

	// Initialize MongoDB
	db := mongodb.NewMongoDB(os.Getenv("DB_Name"))

	middleware := core_util.NewMiddlewareUtil()
	memberUtil := core_util.NewMemberUtil(db)

	externalRoute := gin.NewGinMainRoute(middleware, memberUtil)

	mainRoute := service.NewMainRoute(externalRoute)

	mainRoute.RouteMember()
}

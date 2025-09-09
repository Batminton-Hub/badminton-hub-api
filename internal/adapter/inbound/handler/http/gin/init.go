package gin

import (
	"Badminton-Hub/internal/core/port"

	"github.com/gin-gonic/gin"
)

type MainRoute struct {
	engine         *gin.Engine
	authentication AuthenticationSystemController
	redirect       RedirectController
}

func NewGinRoute(
	authenticationSystem port.AuthenticationSystem,
	redirect port.RedirectService,
) *MainRoute {
	return &MainRoute{
		authentication: &AuthenticationSystem{authenticationSystem},
		redirect:       &Redirect{redirect},
	}
}

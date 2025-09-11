package gin

import (
	"Badminton-Hub/internal/core/port"

	"github.com/gin-gonic/gin"
)

type MainRoute struct {
	engine         *gin.Engine
	authentication AuthenticationSystemController
	redirect       RedirectController
	member         MemberController
}

func NewGinRoute(
	authenticationSystem port.AuthenticationSystem,
	redirect port.RedirectService,
	member port.MemberService,
) *MainRoute {
	response := &MainRoute{
		authentication: &AuthenticationSystem{authenticationSystem},
		redirect:       &Redirect{redirect},
		member:         &Member{member},
	}
	return response
}

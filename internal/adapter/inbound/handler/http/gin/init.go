package gin

import (
	"Badminton-Hub/internal/core/port"
)

type MainRoute struct {
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

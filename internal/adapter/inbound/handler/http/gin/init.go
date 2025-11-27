package gin

import (
	"Badminton-Hub/internal/core/port"
)

type MainRoute struct {
	healthCheck    HealthCheckController
	authentication AuthenticationSystemController
	redirect       RedirectController
	member         MemberController
	observability  ObservabilityController
	gang           GangController
}

func NewGinRoute(
	authenticationSystem port.AuthenticationSystem,
	redirect port.RedirectService,
	member port.MemberService,
	observability port.Observability,
	gangService port.GangService,
) *MainRoute {
	response := &MainRoute{
		healthCheck:    &HealthCheck{observability},
		authentication: &AuthenticationSystem{authenticationSystem, observability},
		redirect:       &Redirect{redirect},
		member:         &Member{member},
		observability:  &Observability{observability},
		gang:           &Gang{observability, gangService},
	}
	return response
}

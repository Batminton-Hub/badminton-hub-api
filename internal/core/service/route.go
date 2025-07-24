package service

import "Badminton-Hub/internal/core/port"

type MainRoute struct {
	ExternalRoute port.MainRoute
}

func NewMainRoute(
	externalRoute port.MainRoute,
) *MainRoute {
	return &MainRoute{
		ExternalRoute: externalRoute,
	}
}

func (m *MainRoute) RouteMember() {
	m.ExternalRoute.RouteMember()
}

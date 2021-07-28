package bootstrap

import (
	"github.com/hywmongous/example-service/internal/application/routes"
	"go.uber.org/fx"
)

var RouteOptions = fx.Options(
	fx.Provide(RoutesFactory),
	fx.Provide(routes.AccountRoutesFactory),
	fx.Provide(routes.AuthenticationRoutesFactory),
	fx.Provide(routes.SessionRoutesFactory),
	fx.Provide(routes.TicketRoutesFactory),
)

type Route interface {
	Setup()
}

type Routes []Route

func RoutesFactory(
	accountRoutes routes.AccountRoutes,
	authenticationRoutes routes.AuthenticationRoutes,
	sessionRoutes routes.SessionRoutes,
	ticketRoutes routes.TicketRoutes,
) Routes {
	return Routes{
		accountRoutes,
		authenticationRoutes,
		sessionRoutes,
		ticketRoutes,
	}
}

func (
	routes Routes,
) Setup() {
	for _, route := range routes {
		route.Setup()
	}
}

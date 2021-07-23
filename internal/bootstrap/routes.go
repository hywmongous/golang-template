package bootstrap

import (
	application "github.com/hywmongous/example-service/internal/identity/application/routes"
	"go.uber.org/fx"
)

var RouteOptions = fx.Options(
	fx.Provide(RoutesFactory),
	fx.Provide(application.AccountRoutesFactory),
	fx.Provide(application.AuthenticationRoutesFactory),
	fx.Provide(application.SessionRoutesFactory),
	fx.Provide(application.TicketRoutesFactory),
)

type Route interface {
	Setup()
}

type Routes []Route

func RoutesFactory(
	accountRoutes application.AccountRoutes,
	authenticationRoutes application.AuthenticationRoutes,
	sessionRoutes application.SessionRoutes,
	ticketRoutes application.TicketRoutes,
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

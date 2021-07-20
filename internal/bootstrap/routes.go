package bootstrap

import (
	application "github.com/hywmongous/example-service/internal/identity/application/routes"
	"go.uber.org/fx"
)

var RouteOptions = fx.Options(
	fx.Provide(RoutesFactory),
	fx.Provide(application.AccountRoutesFactory),
	fx.Provide(application.SessionRoutesFactory),
	fx.Provide(application.TicketRoutesFactory),
)

type Route interface {
	Setup()
}

type Routes []Route

func RoutesFactory(
	accountRoutes application.AccountRoutes,
	authenticationRoutes application.SessionRoutes,
) Routes {
	return Routes{
		accountRoutes,
		authenticationRoutes,
	}
}

func (
	routes Routes,
) Setup() {
	for _, route := range routes {
		route.Setup()
	}
}

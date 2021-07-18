package bootstrap

import (
	identity "github.com/hywmongous/example-service/internal/identity/application/routes"
	"go.uber.org/fx"
)

var RouteOptions = fx.Options(
	fx.Provide(identity.AccountRoutesFactory),
	fx.Provide(identity.AuthenticationRoutesFactory),
	fx.Provide(RoutesFactory),
)

type Route interface {
	Setup()
}

type Routes []Route

func RoutesFactory(
	accountRoutes identity.AccountRoutes,
	authenticationRoutes identity.AuthenticationRoutes,
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

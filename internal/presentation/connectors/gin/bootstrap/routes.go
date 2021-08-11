package bootstrap

import (
	"github.com/hywmongous/example-service/internal/presentation/connectors/gin/routes"
	"go.uber.org/fx"
)

var RouteOptions = fx.Options(
	fx.Provide(RoutesFactory),
	fx.Provide(routes.CreateAccountRoutes),
	fx.Provide(routes.CreateAuthenticationRoutes),
	fx.Provide(routes.CreateSessionRoutes),
	fx.Provide(routes.CreateTicketRoutes),
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

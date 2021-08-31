package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hywmongous/example-service/internal/infrastructure/services"
	"github.com/hywmongous/example-service/internal/presentation/gin/controllers"
	"github.com/hywmongous/example-service/internal/presentation/gin/routes"
	"go.uber.org/fx"
)

func Run() {
	fx.New(Module).Run()
}

var Module = fx.Options(
	ControllerOptions,
	RouteOptions,
	InfrastructureOptions,
	engineOptions,
	fx.Invoke(bootstrap),
)

var engineOptions = fx.Option(
	fx.Provide(gin.New),
)

var InfrastructureOptions = fx.Options(
	fx.Provide(services.JWTServiceFactory),
)

var ControllerOptions = fx.Options(
	fx.Provide(controllers.AccountControllerFactory),
	fx.Provide(controllers.AuthenticationControllerFactory),
	fx.Provide(controllers.SessionControllerFactory),
	fx.Provide(controllers.TicketControllerFactory),
)

var RouteOptions = fx.Options(
	fx.Provide(routes.RoutesFactory),
	fx.Provide(routes.CreateAccountRoutes),
	fx.Provide(routes.CreateAuthenticationRoutes),
	fx.Provide(routes.CreateSessionRoutes),
	fx.Provide(routes.CreateTicketRoutes),
)

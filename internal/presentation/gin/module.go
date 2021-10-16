package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hywmongous/example-service/internal/application"
	"github.com/hywmongous/example-service/internal/infrastructure"
	"github.com/hywmongous/example-service/internal/infrastructure/cqrs"
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
	ActorOptions,
	fx.Invoke(bootstrap),
)

var engineOptions = fx.Option(
	fx.Provide(gin.New),
)

var ActorOptions = fx.Options(
	fx.Provide(application.UnregisteredUserFactory),
	fx.Provide(application.RegisteredUserFactory),
)

var InfrastructureOptions = fx.Options(
	fx.Provide(services.JWTServiceFactory),
	fx.Provide(infrastructure.KafkaStreamFactory),
	fx.Provide(infrastructure.MongoStoreFactory),
	fx.Provide(cqrs.IdentityRepositoryFactory),
	fx.Provide(infrastructure.UnitOfWorkFactory),
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

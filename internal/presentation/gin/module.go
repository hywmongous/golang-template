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
	var engineOptions = fx.Option(
		fx.Provide(gin.New),
	)

	var actorOptions = fx.Options(
		fx.Provide(application.UnregisteredUserFactory),
		fx.Provide(application.RegisteredUserFactory),
	)

	var infrastructureOptions = fx.Options(
		fx.Provide(services.JWTServiceFactory),
		fx.Provide(infrastructure.KafkaStreamFactory),
		fx.Provide(infrastructure.MongoStoreFactory),
		fx.Provide(cqrs.IdentityRepositoryFactory),
		fx.Provide(infrastructure.UnitOfWorkFactory),
	)

	var controllerOptions = fx.Options(
		fx.Provide(controllers.AccountControllerFactory),
		fx.Provide(controllers.AuthenticationControllerFactory),
		fx.Provide(controllers.SessionControllerFactory),
		fx.Provide(controllers.TicketControllerFactory),
	)

	var routeOptions = fx.Options(
		fx.Provide(routes.RoutesFactory),
		fx.Provide(routes.CreateAccountRoutes),
		fx.Provide(routes.CreateAuthenticationRoutes),
		fx.Provide(routes.CreateSessionRoutes),
		fx.Provide(routes.CreateTicketRoutes),
	)

	var module = fx.Options(
		controllerOptions,
		routeOptions,
		infrastructureOptions,
		engineOptions,
		actorOptions,
		fx.Invoke(bootstrap),
	)

	fx.New(module).Run()
}

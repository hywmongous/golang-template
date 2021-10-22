package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hywmongous/example-service/internal/application"
	"github.com/hywmongous/example-service/internal/infrastructure"
	"github.com/hywmongous/example-service/internal/infrastructure/services"
	"github.com/hywmongous/example-service/internal/presentation/gin/controllers"
	"github.com/hywmongous/example-service/internal/presentation/gin/routes"
	"go.uber.org/fx"
)

func Run() {
	engineOptions := fx.Provide(gin.New)

	actorOptions := fx.Options(
		fx.Provide(application.UnregisteredUserFactory),
		fx.Provide(application.RegisteredUserFactory),
	)

	infrastructureOptions := fx.Options(
		fx.Provide(services.JWTServiceFactory),
		// fx.Provide(mediator.Create, infrastructure.UnitOfWorkFactory, cqrs.IdentityRepositoryFactory),
		fx.Provide(infrastructure.UnitOfWorkFactory),
		fx.Provide(infrastructure.KafkaStreamFactory),
		fx.Provide(infrastructure.MongoStoreFactory),
	)

	controllerOptions := fx.Options(
		fx.Provide(controllers.AccountControllerFactory),
		fx.Provide(controllers.AuthenticationControllerFactory),
		fx.Provide(controllers.SessionControllerFactory),
		fx.Provide(controllers.TicketControllerFactory),
	)

	routeOptions := fx.Options(
		fx.Provide(routes.Factory),
		fx.Provide(routes.CreateAccountRoutes),
		fx.Provide(routes.CreateAuthenticationRoutes),
		fx.Provide(routes.CreateSessionRoutes),
		fx.Provide(routes.CreateTicketRoutes),
	)

	module := fx.Options(
		controllerOptions,
		routeOptions,
		infrastructureOptions,
		engineOptions,
		actorOptions,
		fx.Invoke(bootstrap),
	)

	fx.New(module).Run()
}

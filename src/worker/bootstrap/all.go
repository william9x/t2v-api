package bootstrap

import (
	"context"
	adapter "github.com/Braly-Ltd/t2v-api-adapter"
	"github.com/Braly-Ltd/t2v-api-adapter/clients"
	"github.com/Braly-Ltd/t2v-api-adapter/firebase"
	adapterProps "github.com/Braly-Ltd/t2v-api-adapter/properties"
	"github.com/Braly-Ltd/t2v-api-core/ports"
	"github.com/Braly-Ltd/t2v-api-worker/handlers"
	"github.com/Braly-Ltd/t2v-api-worker/properties"
	"github.com/Braly-Ltd/t2v-api-worker/routers"
	"github.com/golibs-starter/golib"
	golibgin "github.com/golibs-starter/golib-gin"
	"github.com/golibs-starter/golib/log"
	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
)

func All() fx.Option {
	return fx.Options(
		golib.AppOpt(),
		golib.PropertiesOpt(),
		golib.LoggingOpt(),
		golib.EventOpt(),
		golib.BuildInfoOpt(Version, CommitHash, BuildTime),
		golib.ActuatorEndpointOpt(),
		golib.HttpRequestLogOpt(),
		AsynqWorkerOpt(),

		// Provide all application properties
		golib.ProvideProps(properties.NewFileProperties),
		golib.ProvideProps(properties.NewWorkerProperties),
		golib.ProvideProps(adapterProps.NewMinIOProperties),
		golib.ProvideProps(adapterProps.NewAsynqProperties),
		golib.ProvideProps(adapterProps.NewAnimateLCMProperties),
		golib.ProvideProps(adapterProps.NewFirebaseProperties),
		golib.ProvideProps(adapterProps.NewMongoProperties),

		// Provide clients
		fx.Provide(clients.NewMinIOClient),
		fx.Provide(clients.NewAsynqClient),
		fx.Provide(clients.NewHTTPClient),
		fx.Provide(clients.NewMongoClient),

		fx.Provide(firebase.NewFirebaseApplication),
		fx.Provide(firebase.NewFirebaseMessagingClient),

		// Provide port's implements
		fx.Provide(fx.Annotate(
			adapter.NewMinIOAdapter, fx.As(new(ports.ObjectStoragePort))),
		),
		fx.Provide(fx.Annotate(
			adapter.NewAsynqAdapter, fx.As(new(ports.TaskQueuePort))),
		),
		fx.Provide(fx.Annotate(
			adapter.NewAnimateLCMAdapter, fx.As(new(ports.InferencePort))),
		),
		fx.Provide(fx.Annotate(
			firebase.NewFirebaseMessagingClient, fx.As(new(ports.NotificationPort))),
		),
		fx.Provide(fx.Annotate(
			adapter.NewNotificationSubscriptionAdapter, fx.As(new(ports.NotificationSubscriptionPort))),
		),

		// Provide task handlers
		handlers.ProvideHandler(handlers.NewT2VHandler),

		ProvideAsynqWorker(),

		// Provide use cases

		// Provide controllers, these controllers will be used
		// when register router was invoked

		// Provide gin http server auto config,
		// actuator endpoints and application routers
		golibgin.GinHttpServerOpt(),
		fx.Invoke(routers.RegisterGinRouters),

		// Graceful shutdown.
		// OnStop hooks will run in reverse order.
		OnStopAsynqWorker(),
		golibgin.OnStopHttpServerOpt(),
		fx.Invoke(OnStopMongoHook),
		fx.Invoke(OnStopRedisHook),
	)
}

func OnStopMongoHook(lc fx.Lifecycle, client *mongo.Client) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Infof("Disconnecting mongo client")
			if err := client.Disconnect(ctx); err != nil {
				log.Errorf("Could not disconnect mongo client, error [%v]", err)
			}
			return nil
		},
	})
}

func OnStopRedisHook(lc fx.Lifecycle, client *asynq.Client) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Infof("Disconnecting asynq client")
			if err := client.Close(); err != nil {
				log.Errorf("Could not disconnect asynq client, error [%v]", err)
			}
			return nil
		},
	})
}

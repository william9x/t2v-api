package bootstrap

import (
	"context"
	"errors"
	adapter "github.com/Braly-Ltd/t2v-api-adapter"
	"github.com/Braly-Ltd/t2v-api-adapter/clients"
	"github.com/Braly-Ltd/t2v-api-adapter/firebase"
	adapterProps "github.com/Braly-Ltd/t2v-api-adapter/properties"
	"github.com/Braly-Ltd/t2v-api-core/ports"
	"github.com/Braly-Ltd/t2v-api-public/controllers"
	"github.com/Braly-Ltd/t2v-api-public/middlewares"
	"github.com/Braly-Ltd/t2v-api-public/properties"
	"github.com/Braly-Ltd/t2v-api-public/routers"
	"github.com/Braly-Ltd/t2v-api-public/services"
	"github.com/Braly-Ltd/t2v-api-public/validators"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib"
	golibgin "github.com/golibs-starter/golib-gin"
	"github.com/golibs-starter/golib/log"
	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"net/http"
	"os"
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

		// Provide all application properties
		golib.ProvideProps(properties.NewSwaggerProperties),
		golib.ProvideProps(properties.NewTLSProperties),
		golib.ProvideProps(properties.NewModelProperties),
		golib.ProvideProps(properties.NewPromptProperties),
		golib.ProvideProps(properties.NewMiddlewaresProperties),
		golib.ProvideProps(adapterProps.NewMinIOProperties),
		golib.ProvideProps(adapterProps.NewAsynqProperties),
		golib.ProvideProps(adapterProps.NewFirebaseProperties),
		golib.ProvideProps(adapterProps.NewMongoProperties),

		// Provide clients
		fx.Provide(clients.NewHTTPClient),
		fx.Provide(clients.NewMinIOClient),
		fx.Provide(clients.NewAsynqClient),
		fx.Provide(clients.NewAsynqInspector),
		fx.Provide(clients.NewMongoClient),
		fx.Provide(firebase.NewFirebaseApplication),
		fx.Provide(firebase.NewFirebaseAuthClient),
		fx.Provide(firebase.NewFirebaseMessagingClient),

		// Provide port's implements
		fx.Provide(fx.Annotate(
			adapter.NewMinIOAdapter, fx.As(new(ports.ObjectStoragePort))),
		),
		fx.Provide(fx.Annotate(
			adapter.NewAsynqAdapter, fx.As(new(ports.TaskQueuePort))),
		),
		fx.Provide(fx.Annotate(
			firebase.NewFirebaseAuthClient, fx.As(new(ports.AuthenticationPort))),
		),
		fx.Provide(fx.Annotate(
			adapter.NewNotificationSubscriptionAdapter, fx.As(new(ports.NotificationSubscriptionPort))),
		),
		fx.Provide(fx.Annotate(
			adapter.NewTaskInfoAdapter, fx.As(new(ports.TaskInfoRepositoryPort))),
		),

		// Provide use cases
		fx.Provide(services.NewInferenceService),
		fx.Provide(services.NewProfanityDetector),

		// Provide controllers, these controllers will be used
		// when register router was invoked
		fx.Provide(controllers.NewInferenceController),
		fx.Provide(controllers.NewModelController),
		fx.Provide(controllers.NewPromptController),
		fx.Provide(controllers.NewNotificationController),
		fx.Provide(controllers.NewNotificationControllerV2),

		// Provide gin http server auto config,
		// actuator endpoints and application routers
		GinHttpServerOpt(),
		fx.Invoke(routers.RegisterGinRouters),

		// Register custom validators
		fx.Invoke(validators.RegisterFormValidators),

		// Graceful shutdown.
		// OnStop hooks will run in reverse order.
		golibgin.OnStopHttpServerOpt(),
		fx.Invoke(OnStopMongoHook),
		fx.Invoke(OnStopRedisHook),
	)
}

func GinHttpServerOpt() fx.Option {
	return fx.Options(
		fx.Provide(golibgin.NewGinEngine),
		fx.Provide(golibgin.NewHTTPServer),
		fx.Invoke(RegisterMiddlewares),
		fx.Invoke(golibgin.RegisterHandlers),
		fx.Invoke(OnStartHttpsServerHook),
	)
}

func OnStartHttpsServerHook(lc fx.Lifecycle, app *golib.App, httpServer *http.Server, tls *properties.TLSProperties) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Infof("Application will be served at %s. Service name: %s, service path: %s",
				httpServer.Addr, app.Name(), app.Path())
			go func() {
				if tls.Enabled {
					log.Infof("Activating Release mode and TLS")
					gin.SetMode(gin.ReleaseMode)
					if _, err := os.Stat("/app/ssl/fullchain.pem"); errors.Is(err, os.ErrNotExist) {
						log.Fatalf("fullchain not exist")
					}
					if _, err := os.Stat("/app/ssl/privkey.pem"); errors.Is(err, os.ErrNotExist) {
						log.Fatalf("privkey not exist")
					}
					if err := httpServer.ListenAndServeTLS(tls.CertFile, tls.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
						log.Errorf("Could not serve HTTP request at %s, error [%v]", httpServer.Addr, err)
					}
				} else {
					if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
						log.Errorf("Could not serve HTTP request at %s, error [%v]", httpServer.Addr, err)
					}
				}
				log.Infof("Stopped HTTP Server %s", httpServer.Addr)
			}()
			return nil
		},
	})
}

func RegisterMiddlewares(app *golib.App) {
	app.AddHandler(
		middlewares.AddCustomHeaders(),
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

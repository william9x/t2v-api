package routers

import (
	"github.com/Braly-Ltd/t2v-api-core/ports"
	"github.com/Braly-Ltd/t2v-api-public/controllers"
	"github.com/Braly-Ltd/t2v-api-public/docs"
	"github.com/Braly-Ltd/t2v-api-public/middlewares"
	"github.com/Braly-Ltd/t2v-api-public/properties"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib"
	"github.com/golibs-starter/golib/web/actuator"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

type RegisterRoutersIn struct {
	fx.In
	App                   *golib.App
	Engine                *gin.Engine
	SwaggerProps          *properties.SwaggerProperties
	MiddlewaresProperties *properties.MiddlewaresProperties

	Actuator            *actuator.Endpoint
	InferenceController *controllers.InferenceController
	ModelController     *controllers.ModelController
	PromptController    *controllers.PromptController

	AuthenticationPort ports.AuthenticationPort
}

func RegisterGinRouters(p RegisterRoutersIn) {
	group := p.Engine.Group(p.App.Path())
	group.GET("/actuator/health", gin.WrapF(p.Actuator.Health))
	group.GET("/actuator/info", gin.WrapF(p.Actuator.Info))

	if p.SwaggerProps.Enabled {
		docs.SwaggerInfo.BasePath = p.App.Path()
		group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	apiV1Group := group.Group("/api/v1")
	apiV1Group.Use(middlewares.Authenticate(p.AuthenticationPort, p.MiddlewaresProperties))
	apiV1Group.Use(middlewares.LogExtraHeader())

	// Model APIs
	apiV1Group.GET("/models", p.ModelController.GetModels)

	// Inference APIs
	apiV1Group.GET("/infer/:id", p.InferenceController.GetInference)
	apiV1Group.GET("/infer", p.InferenceController.FilterInference)
	apiV1Group.POST("/infer", p.InferenceController.CreateInference)

	// Prompt APIs
	apiV1Group.GET("/prompts/suggest", p.PromptController.GetRandomPrompt)
	apiV1Group.GET("/prompts/profanity", p.PromptController.CheckProfanity)
}

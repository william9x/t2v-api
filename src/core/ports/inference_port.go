package ports

import (
	"context"
	"github.com/Braly-Ltd/t2v-api-core/entities"
)

type InferencePort interface {
	Infer(context.Context, entities.InferenceCommand) (entities.InferenceResult, error)
}

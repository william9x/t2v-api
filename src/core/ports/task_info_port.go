package ports

import (
	"context"
	"github.com/Braly-Ltd/t2v-api-core/entities"
)

type TaskInfoRepositoryPort interface {
	Save(ctx context.Context, taskInfo entities.InferTaskInfo) error
}

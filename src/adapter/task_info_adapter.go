package adapter

import (
	"context"
	"errors"
	"fmt"
	"github.com/Braly-Ltd/t2v-api-core/entities"
	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TaskInfoAdapter ...
type TaskInfoAdapter struct {
	client    *mongo.Client
	inspector *asynq.Inspector
}

// NewTaskInfoAdapter ...
func NewTaskInfoAdapter(client *mongo.Client, inspector *asynq.Inspector) *TaskInfoAdapter {
	return &TaskInfoAdapter{
		client:    client,
		inspector: inspector,
	}
}

func (r *TaskInfoAdapter) FindByIDAndQueue(ctx context.Context, queue, id string) (entities.InferTaskInfo, error) {
	taskInfo, err := r.inspector.GetTaskInfo(queue, id)

	if errors.Is(err, asynq.ErrQueueNotFound) {
		return entities.InferTaskInfo{}, fmt.Errorf("queue %q not found", queue)
	}

	if errors.Is(err, asynq.ErrTaskNotFound) {
		result := r.collection().FindOne(ctx, bson.M{"_id": id, "queue": queue})
		if err := result.Err(); err != nil {
			return entities.InferTaskInfo{}, fmt.Errorf("task %q not found in queue %q", id, queue)
		}

		var taskInfo entities.InferTaskInfo
		if err := result.Decode(&taskInfo); err != nil {
			return entities.InferTaskInfo{}, fmt.Errorf("failed to decode task info: %w", err)
		}
		return taskInfo, nil
	}

	if err != nil {
		return entities.InferTaskInfo{}, err
	}

	return entities.InferTaskInfo{
		TaskID:   taskInfo.ID,
		Queue:    taskInfo.Queue,
		Type:     taskInfo.Type,
		MaxRetry: taskInfo.MaxRetry,
		//Deadline:      taskInfo.Deadline,
		//RetentionUtil: taskInfo.Retention,
		Status: taskInfo.State.String(),
		//EnqueuedAt:    taskInfo.EnqueuedAt,
		CompletedAt: taskInfo.CompletedAt.UnixMilli(),
		Payload:     entities.InferPayload{},
	}, nil
}

// Save ...
func (r *TaskInfoAdapter) Save(ctx context.Context, taskInfo entities.InferTaskInfo) error {
	_, err := r.collection().InsertOne(ctx, taskInfo)
	return err
}

func (r *TaskInfoAdapter) collection() *mongo.Collection {
	return r.client.Database(mongodbDatabaseName).Collection(mongodbTaskInfoCollection)
}

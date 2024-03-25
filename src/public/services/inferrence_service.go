package services

import (
	"context"
	"fmt"
	"github.com/Braly-Ltd/t2v-api-core/constants"
	"github.com/Braly-Ltd/t2v-api-core/entities"
	"github.com/Braly-Ltd/t2v-api-core/ports"
	"github.com/Braly-Ltd/t2v-api-core/utils"
	"github.com/Braly-Ltd/t2v-api-public/requests"
	"github.com/Braly-Ltd/t2v-api-public/resources"
	"github.com/hibiken/asynq"
	"time"
)

type InferenceService struct {
	objectStoragePort ports.ObjectStoragePort
	taskQueuePort     ports.TaskQueuePort
}

func NewInferenceService(
	objectStoragePort ports.ObjectStoragePort,
	taskQueuePort ports.TaskQueuePort,
) *InferenceService {
	return &InferenceService{
		objectStoragePort: objectStoragePort,
		taskQueuePort:     taskQueuePort,
	}
}

// GetInference ...
func (r *InferenceService) GetInference(ctx context.Context, queueId, id string) (*asynq.TaskInfo, error) {
	return r.taskQueuePort.GetTask(ctx, queueId, id)
}

// FilterInference ...
func (r *InferenceService) FilterInference(ctx context.Context, ids []string) ([]*asynq.TaskInfo, error) {
	tasks := make([]*asynq.TaskInfo, 0, len(ids))
	for _, id := range ids {
		queueId, inferId := utils.ExtractInferenceKey(id)
		task, err := r.taskQueuePort.GetTask(ctx, queueId, inferId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// CreateInference ...
func (r *InferenceService) CreateInference(
	ctx context.Context,
	req requests.CreateInferenceRequest,
) (resources.CreateInference, error) {
	taskId, err := utils.NewUUID()
	if err != nil {
		return resources.CreateInference{}, fmt.Errorf("generate task id error: %v", err)
	}

	targetFileName := fmt.Sprintf("target/%s%s", taskId, constants.FileTypeDefault)
	targetFileURL, err := r.objectStoragePort.GetPreSignedObject(ctx, targetFileName)
	if err != nil {
		return resources.CreateInference{}, fmt.Errorf("get pre-signed target object error: %v", err)
	}

	packed, err := newPackedInferPayload(req, targetFileName, targetFileURL)
	if err != nil {
		return resources.CreateInference{}, fmt.Errorf("pack payload error: %v", err)
	}

	// TODO: Select queue based on model type
	queue := string(constants.QueueTypeDefault)
	maxRetry := 0
	deadline := time.Now().Add(10 * time.Minute)
	retention := 24 * time.Hour

	taskOpts := []asynq.Option{
		asynq.TaskID(taskId),
		asynq.Queue(queue),
		asynq.MaxRetry(maxRetry),
		asynq.Deadline(deadline),
		asynq.Retention(retention),
	}
	task := asynq.NewTask(req.Type, packed, taskOpts...)
	if err := r.taskQueuePort.Enqueue(ctx, task); err != nil {
		return resources.CreateInference{}, err
	}

	return resources.CreateInference{
		ID:        utils.BuildInferenceKey(queue, taskId),
		Model:     req.Model,
		Type:      req.Type,
		Status:    asynq.TaskStatePending.String(),
		MaxRetry:  maxRetry,
		Deadline:  deadline.Format(time.RFC3339),
		Retention: time.Now().Add(retention).Format(time.RFC3339),
	}, nil
}

func newPackedInferPayload(
	req requests.CreateInferenceRequest,
	targetFileName string,
	targetFileURL string,
) ([]byte, error) {
	payload := newInferPayload(req, targetFileName, targetFileURL)
	return payload.Packed()
}

func newInferPayload(
	req requests.CreateInferenceRequest,
	targetFileName string,
	targetFileURL string,
) *entities.InferPayload {
	return &entities.InferPayload{
		Model:             req.Model,
		Prompt:            req.Prompt,
		NegativePrompt:    req.NegativePrompt,
		NumInferenceSteps: req.NumInferenceSteps,
		Height:            req.Height,
		Width:             req.Width,
		NumFrames:         req.NumFrames,
		GuidanceScale:     req.GuidanceScale,
		TargetFileName:    targetFileName,
		TargetFileURL:     targetFileURL,
		EnqueuedAt:        time.Now().UnixMilli(),
		UserID:            req.UserID,
		Agent:             req.Agent,
	}
}

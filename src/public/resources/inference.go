package resources

import (
	"github.com/Braly-Ltd/t2v-api-core/entities"
	"github.com/Braly-Ltd/t2v-api-core/utils"
	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

// CreateInference ...
type CreateInference struct {
	ID        string `json:"id,omitempty"`
	Model     string `json:"model,omitempty"`
	Type      string `json:"type,omitempty"`
	Status    string `json:"status,omitempty"` // Status of the task. Values: active, pending, scheduled, retry, archived, completed
	MaxRetry  int    `json:"max_retry"`
	Deadline  string `json:"deadline,omitempty"` // Deadline for completing the task
	Retention string `json:"retention"`          // Retention in hours for how long to store the task info
}

// Inference ...
type Inference struct {
	ID            string `json:"id,omitempty"`
	Model         string `json:"model,omitempty"`
	Type          string `json:"type,omitempty"`
	Status        string `json:"status,omitempty"` // Status of the task. Values: active, pending, scheduled, retry, archived, completed
	MaxRetry      int    `json:"max_retry"`
	Deadline      string `json:"deadline,omitempty"` // Deadline for completing the task
	Retried       int    `json:"retried"`
	LastErr       string `json:"last_err,omitempty"`
	LastFailedAt  string `json:"last_failed_at,omitempty"`
	TargetFileURL string `json:"target_file_url,omitempty"`
	EnqueuedAt    string `json:"enqueued_at,omitempty"`
	CompletedAt   string `json:"completed_at,omitempty"`
}

func NewFromTaskInfoList(infoList []*asynq.TaskInfo) ([]*Inference, error) {
	inferences := make([]*Inference, 0, len(infoList))
	for _, item := range infoList {
		inference, err := NewFromTaskInfo(item)
		if err != nil {
			return nil, err
		}
		inferences = append(inferences, inference)
	}
	return inferences, nil
}

func NewFromTaskInfo(info *asynq.TaskInfo) (*Inference, error) {
	var payload entities.InferPayload
	if err := msgpack.Unmarshal(info.Payload, &payload); err != nil {
		return nil, err
	}

	var failedAt time.Time
	if info.LastFailedAt.UnixMilli() > 0 {
		failedAt = time.UnixMilli(info.LastFailedAt.UnixMilli())
	}

	return &Inference{
		ID:            utils.BuildInferenceKey(info.Queue, info.ID),
		Model:         payload.Model,
		Type:          info.Type,
		Status:        info.State.String(),
		MaxRetry:      info.MaxRetry,
		Deadline:      info.Deadline.Format(time.RFC3339),
		Retried:       info.Retried,
		LastErr:       info.LastErr,
		LastFailedAt:  failedAt.Format(time.RFC3339),
		TargetFileURL: payload.TargetFileURL,
		EnqueuedAt:    time.UnixMilli(payload.EnqueuedAt).Format(time.RFC3339),
		CompletedAt:   info.CompletedAt.Format(time.RFC3339),
	}, nil
}

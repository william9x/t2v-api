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
	LastFailedAt  int64  `json:"last_failed_at,omitempty"`
	TargetFileURL string `json:"target_file_url,omitempty"`
}

func NewFromTaskInfo(info *asynq.TaskInfo) (*Inference, error) {
	var payload entities.InferPayload
	if err := msgpack.Unmarshal(info.Payload, &payload); err != nil {
		return nil, err
	}

	var failedAt int64 = 0
	if info.LastFailedAt.UnixMilli() > 0 {
		failedAt = info.LastFailedAt.UnixMilli()
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
		LastFailedAt:  failedAt,
		TargetFileURL: payload.TargetFileURL,
	}, nil
}

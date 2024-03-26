package entities

type InferTaskInfo struct {
	TaskID        string       `bson:"_id"`
	Queue         string       `bson:"queue"`
	Type          string       `bson:"type"`
	MaxRetry      int          `bson:"max_retry"`
	Deadline      string       `bson:"deadline"`
	RetentionUtil string       `bson:"retention"`
	Status        string       `bson:"status"`
	EnqueuedAt    int64        `bson:"enqueued_at"`
	CompletedAt   int64        `bson:"completed_at"`
	Payload       InferPayload `bson:"payload"`
}

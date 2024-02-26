package constants

type TaskType string

const (
	TaskTypeT2V     TaskType = "t2v"
	TaskTypeI2V     TaskType = "i2v"
	TaskTypeV2V     TaskType = "v2v"
	TaskTypeUpscale TaskType = "upscale"
)

const (
	TaskIDSeparator = ":"
)

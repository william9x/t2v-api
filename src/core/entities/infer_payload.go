package entities

import (
	"github.com/vmihailenco/msgpack/v5"
)

type InferPayload struct {
	UserID            string  `bson:"user_id"`
	Agent             string  `bson:"agent"`
	Model             string  `bson:"model"`
	Prompt            string  `bson:"prompt"`
	NegativePrompt    string  `bson:"negative_prompt"`
	NumInferenceSteps int     `bson:"num_inference_steps"`
	NumFrames         int     `bson:"num_frames"`
	Width             int     `bson:"width"`
	Height            int     `bson:"height"`
	GuidanceScale     float32 `bson:"guidance_scale"`
	TargetFileName    string  `bson:"target_file_name"`
	TargetFileURL     string  `bson:"target_file_url"`

	// @Deprecated
	EnqueuedAt int64 `bson:"-"`
}

func (p *InferPayload) Packed() ([]byte, error) {
	return msgpack.Marshal(p)
}

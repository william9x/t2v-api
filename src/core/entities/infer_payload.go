package entities

import (
	"github.com/vmihailenco/msgpack/v5"
)

type InferPayload struct {
	UserID            string
	Model             string
	Prompt            string
	NegativePrompt    string
	NumInferenceSteps int
	NumFrames         int
	Width             int
	Height            int
	GuidanceScale     float32
	TargetFileName    string
	TargetFileURL     string
	EnqueuedAt        int64
}

func (p *InferPayload) Packed() ([]byte, error) {
	return msgpack.Marshal(p)
}

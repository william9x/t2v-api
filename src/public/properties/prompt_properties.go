package properties

import (
	"github.com/golibs-starter/golib/config"
)

func NewPromptProperties(loader config.Loader) (*PromptProperties, error) {
	props := PromptProperties{}
	err := loader.Bind(&props)
	return &props, err
}

type PromptProperties struct {
	Data []*PromptData
}

type PromptData struct {
	Prompt            string
	NegPrompt         string
	ModelID           string
	NumInferenceSteps int
	NumFrames         int
	Width             int
	Height            int
	GuidanceScale     float32
}

func (t *PromptProperties) Prefix() string {
	return "app.prompts"
}

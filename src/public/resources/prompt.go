package resources

import "github.com/Braly-Ltd/t2v-api-public/properties"

// Prompt ...
type Prompt struct {
	Prompt            string  `json:"prompt,omitempty"`
	NegPrompt         string  `json:"neg_prompt,omitempty"`
	ModelID           string  `json:"model_id,omitempty"`
	NumInferenceSteps int     `json:"num_inference_steps,omitempty"`
	NumFrames         int     `json:"num_frames,omitempty"`
	Width             int     `json:"width,omitempty"`
	Height            int     `json:"height,omitempty"`
	GuidanceScale     float32 `json:"guidance_scale,omitempty"`
}

func FromPromptProps(props *properties.PromptProperties) []Prompt {
	modelResources := make([]Prompt, len(props.Data))
	for i, prompt := range props.Data {
		modelResources[i] = FromPromptData(prompt)
	}
	return modelResources
}

func FromPromptData(prompt *properties.PromptData) Prompt {
	return Prompt{
		Prompt:            prompt.Prompt,
		NegPrompt:         prompt.NegPrompt,
		ModelID:           prompt.ModelID,
		NumInferenceSteps: prompt.NumInferenceSteps,
		NumFrames:         prompt.NumFrames,
		Width:             prompt.Width,
		Height:            prompt.Height,
		GuidanceScale:     prompt.GuidanceScale,
	}
}

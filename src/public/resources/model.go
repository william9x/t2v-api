package resources

import "github.com/Braly-Ltd/t2v-api-public/properties"

// Model ...
type Model struct {
	ID                string         `json:"id,omitempty"`
	Name              string         `json:"name,omitempty"`
	LogoURL           string         `json:"logo_url,omitempty"`
	NumInferenceSteps int            `json:"num_inference_steps,omitempty"`
	NumFrames         int            `json:"num_frames,omitempty"`
	Width             int            `json:"width,omitempty"`
	Height            int            `json:"height,omitempty"`
	GuidanceScale     float32        `json:"guidance_scale,omitempty"`
	RecomdPrompts     []RecomdPrompt `json:"recomd_prompts,omitempty"`
}

type RecomdPrompt struct {
	Prompt       string `json:"prompt,omitempty"`
	NegPrompt    string `json:"neg_prompt,omitempty"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	AssetURL     string `json:"asset_url,omitempty"`
}

func FromModelProps(props *properties.ModelProperties) []Model {
	modelResources := make([]Model, len(props.Data))
	for i, model := range props.Data {
		recomdPromptResources := make([]RecomdPrompt, len(model.RecomdPrompt))
		for j, prompt := range model.RecomdPrompt {
			recomdPromptResources[j] = RecomdPrompt{
				Prompt:       prompt.Prompt,
				NegPrompt:    prompt.NegPrompt,
				ThumbnailURL: prompt.ThumbnailURL,
				AssetURL:     prompt.AssetURL,
			}
		}

		modelResources[i] = Model{
			ID:                model.ID,
			Name:              model.Name,
			LogoURL:           model.LogoURL,
			NumInferenceSteps: model.NumInferenceSteps,
			NumFrames:         model.NumFrames,
			Width:             model.Width,
			Height:            model.Height,
			GuidanceScale:     model.GuidanceScale,
			RecomdPrompts:     recomdPromptResources,
		}
	}
	return modelResources
}

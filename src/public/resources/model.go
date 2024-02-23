package resources

import "github.com/Braly-Ltd/t2v-api-public/properties"

// Model ...
type Model struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	LogoURL       string         `json:"logo_url,omitempty"`
	RecomdPrompts []RecomdPrompt `json:"recomd_prompts"`
}

type RecomdPrompt struct {
	Prompt    string `json:"prompt,omitempty"`
	NegPrompt string `json:"neg_prompt,omitempty"`
	LogoURL   string `json:"logo_url,omitempty"`
}

func FromModelProps(props *properties.ModelProperties) []Model {
	modelResources := make([]Model, len(props.Data))
	for i, model := range props.Data {
		recomdPromptResources := make([]RecomdPrompt, len(model.RecomdPrompt))
		for j, prompt := range model.RecomdPrompt {
			recomdPromptResources[j] = RecomdPrompt{
				Prompt:    prompt.Prompt,
				NegPrompt: prompt.NegPrompt,
				LogoURL:   prompt.LogoURL,
			}
		}

		modelResources[i] = Model{
			ID:            model.ID,
			Name:          model.Name,
			LogoURL:       model.LogoURL,
			RecomdPrompts: recomdPromptResources,
		}
	}
	return modelResources
}

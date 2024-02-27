package resources

import "github.com/Braly-Ltd/t2v-api-public/properties"

// Model ...
type Model struct {
	ID           string `json:"id,omitempty"`
	TriggerWords string `json:"trigger_words,omitempty"`
}

func FromModelProps(props *properties.ModelProperties) []Model {
	modelResources := make([]Model, len(props.Data))
	for i, model := range props.Data {
		modelResources[i] = Model{
			ID:           model.ID,
			TriggerWords: model.TriggerWords,
		}
	}
	return modelResources
}

package entities

type InferenceCommand struct {
	Prompt            string  `json:"prompt"`
	NegativePrompt    string  `json:"negative_prompt"`
	NumInferenceSteps int     `json:"num_inference_steps"`
	NumFrames         int     `json:"num_frames"`
	Width             int     `json:"width"`
	Height            int     `json:"height"`
	GuidanceScale     float32 `json:"guidance_scale"`
	OutputFilePath    string  `json:"output_file_path"`
	ModelID           string  `json:"model_id"`
}

package requests

// CreateInferenceRequest ...
type CreateInferenceRequest struct {
	Model  string `form:"model,omitempty" binding:"notblank"`
	Type   string `form:"type,omitempty,default=t2v"`
	UserID string `form:"user_id,omitempty"`

	Prompt            string  `form:"prompt,omitempty" binding:"notblank"`
	NegativePrompt    string  `form:"negative_prompt,omitempty"`
	NumInferenceSteps int     `form:"num_inference_steps,omitempty,default=4" binding:"max=200,min=1"`
	NumFrames         int     `form:"num_frames,omitempty,default=16" binding:"max=32,min=16"`
	Width             int     `form:"width,omitempty,default=512"`
	Height            int     `form:"height,omitempty,default=512"`
	GuidanceScale     float32 `form:"guidance_scale,omitempty,default=1.5" binding:"max=100,min=0"`
}

// FilterInferenceRequest ...
type FilterInferenceRequest struct {
	IDs []string `json:"ids,omitempty" binding:"required"`
}

// CheckPromptProfanity ...
type CheckPromptProfanity struct {
	Prompt string `json:"prompt,omitempty" binding:"required"`
}

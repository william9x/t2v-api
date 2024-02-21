package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Braly-Ltd/t2v-api-adapter/properties"
	"github.com/Braly-Ltd/t2v-api-core/entities"
	"net/http"
)

// AnimateLCMAdapter ...
type AnimateLCMAdapter struct {
	client *http.Client
	props  *properties.AnimateLCMProperties
}

// NewAnimateLCMAdapter ...
func NewAnimateLCMAdapter(client *http.Client, props *properties.AnimateLCMProperties) *AnimateLCMAdapter {
	return &AnimateLCMAdapter{client: client, props: props}
}

func (r *AnimateLCMAdapter) Infer(ctx context.Context, payload entities.InferenceCommand) error {
	req := entities.InferenceCommand{
		Prompt:            payload.Prompt,
		NegativePrompt:    payload.NegativePrompt,
		NumInferenceSteps: payload.NumInferenceSteps,
		NumFrames:         payload.NumFrames,
		Width:             payload.Width,
		Height:            payload.Height,
		GuidanceScale:     payload.GuidanceScale,
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return fmt.Errorf("encoding request error: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", r.props.InferURL, buf)
	if err != nil {
		return fmt.Errorf("build http request error: %v", err)
	}

	resp, err := r.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("create inference error: %v", err)
	}

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return fmt.Errorf("create inference failed")
	}

	return nil
}

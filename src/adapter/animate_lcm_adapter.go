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

func (r *AnimateLCMAdapter) Infer(ctx context.Context, cmd entities.InferenceCommand) (entities.InferenceResult, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(cmd); err != nil {
		return entities.InferenceResult{}, fmt.Errorf("encode command failed: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", r.props.InferURL, buf)
	if err != nil {
		return entities.InferenceResult{}, fmt.Errorf("create http request failed: %v", err)
	}

	resp, err := r.client.Do(httpReq)
	if err != nil {
		return entities.InferenceResult{}, fmt.Errorf("do http request failed: %v", err)
	}

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return entities.InferenceResult{}, fmt.Errorf("inferrence error: %d", resp.StatusCode)
	}

	var respData entities.InferenceResult
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return entities.InferenceResult{}, fmt.Errorf("decode response failed: %v", err)
	}

	return respData, nil
}

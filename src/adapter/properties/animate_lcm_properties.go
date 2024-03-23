package properties

import (
	"fmt"
	"github.com/golibs-starter/golib/config"
)

type AnimateLCMProperties struct {
	Endpoint  string
	InferPath string

	InferURL string `default:""`
}

func NewAnimateLCMProperties(loader config.Loader) (*AnimateLCMProperties, error) {
	props := AnimateLCMProperties{}
	err := loader.Bind(&props)
	return &props, err
}

func (r *AnimateLCMProperties) Prefix() string {
	return "app.animatelcm"
}

func (r *AnimateLCMProperties) PostBinding() error {
	r.InferURL = fmt.Sprintf("%s%s", r.Endpoint, r.InferPath)
	return nil
}

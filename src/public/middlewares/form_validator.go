package middlewares

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func RegisterFormValidators() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("notblank", validators.NotBlank); err != nil {
			return err
		}
		if err := v.RegisterValidation("infertype", inferType); err != nil {
			return err
		}
	}
	return nil
}

// inferType ...
func inferType(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return val == "t2v" || val == "i2v" || val == "v2v" || val == "upscale"
}

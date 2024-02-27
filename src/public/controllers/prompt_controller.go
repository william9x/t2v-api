package controllers

import (
	"github.com/Braly-Ltd/t2v-api-public/properties"
	"github.com/Braly-Ltd/t2v-api-public/resources"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/web/response"
	"math/rand"
)

type PromptController struct {
	promptProps *properties.PromptProperties
}

func NewPromptController(
	promptProps *properties.PromptProperties,
) *PromptController {
	return &PromptController{
		promptProps: promptProps,
	}
}

// GetRandomPrompt
//
//	@ID				get-random-prompt
//	@Summary 		Get random prompts
//	@Description
//	@Tags			PromptController
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	response.Response{data=resources.Prompt}
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/prompts/suggest [get]
func (c *PromptController) GetRandomPrompt(ctx *gin.Context) {
	prompts := c.promptProps.Data
	minIdx := 0
	maxIdx := len(prompts)
	idx := rand.Intn(maxIdx-minIdx) + minIdx
	response.Write(ctx.Writer, response.Ok(resources.FromPromptData(prompts[idx])))
}

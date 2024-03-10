package controllers

import (
	"github.com/Braly-Ltd/t2v-api-public/properties"
	"github.com/Braly-Ltd/t2v-api-public/requests"
	"github.com/Braly-Ltd/t2v-api-public/resources"
	goaway "github.com/TwiN/go-away"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/exception"
	"github.com/golibs-starter/golib/web/response"
	"math/rand"
)

type PromptController struct {
	promptProps       *properties.PromptProperties
	profanityDetector *goaway.ProfanityDetector
}

func NewPromptController(
	promptProps *properties.PromptProperties,
	profanityDetector *goaway.ProfanityDetector,
) *PromptController {
	return &PromptController{
		promptProps:       promptProps,
		profanityDetector: profanityDetector,
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

// CheckProfanity
//
//	@ID				check-profanity
//	@Summary 		Check Profanity
//	@Description
//	@Tags			PromptController
//	@Accept			json
//	@Produce		json
//	@Param			req		body    	requests.CheckPromptProfanity     true        "Request body"
//	@Success		200		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/prompts/profanity [get]
func (c *PromptController) CheckProfanity(ctx *gin.Context) {
	var req requests.CheckPromptProfanity
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.WriteError(ctx.Writer, exception.New(40000, err.Error()))
		return
	}

	resp := struct {
		Prompt    string `json:"prompt,omitempty"`
		IsProfane bool   `json:"is_profane"`
	}{
		Prompt:    req.Prompt,
		IsProfane: c.profanityDetector.IsProfane(req.Prompt),
	}
	response.Write(ctx.Writer, response.Ok(resp))
}

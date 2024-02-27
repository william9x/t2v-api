package controllers

import (
	"fmt"
	"github.com/Braly-Ltd/t2v-api-core/utils"
	"github.com/Braly-Ltd/t2v-api-public/properties"
	"github.com/Braly-Ltd/t2v-api-public/requests"
	"github.com/Braly-Ltd/t2v-api-public/resources"
	"github.com/Braly-Ltd/t2v-api-public/services"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/exception"
	"github.com/golibs-starter/golib/log"
	"github.com/golibs-starter/golib/web/response"
)

type InferenceController struct {
	modelProps       *properties.ModelProperties
	inferenceProps   *properties.InferenceProperties
	inferenceService *services.InferenceService
}

func NewInferenceController(
	modelProps *properties.ModelProperties,
	inferenceProps *properties.InferenceProperties,
	inferenceService *services.InferenceService,
) *InferenceController {
	return &InferenceController{
		modelProps:       modelProps,
		inferenceProps:   inferenceProps,
		inferenceService: inferenceService,
	}
}

// GetInference
//
//	@ID				get-inference
//	@Summary 		Get status of an inference task
//	@Description
//	@Tags			InferenceController
//	@Accept			json
//	@Produce		json
//	@Param			id		path    	string     true        "Task ID"
//	@Success		200		{object}	response.Response{data=resources.Inference}
//	@Success		400		{object}	response.Response
//	@Success		404		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/infer/{id} [get]
func (c *InferenceController) GetInference(ctx *gin.Context) {

	queueId, inferId := utils.ExtractInferenceKey(ctx.Param("id"))
	if queueId == "" || inferId == "" {
		response.WriteError(ctx.Writer, exception.New(400, "Invalid infer ID"))
		return
	}

	inferInfo, err := c.inferenceService.GetInference(ctx, queueId, inferId)
	if err != nil {
		response.WriteError(ctx.Writer, exception.New(404, "Task not found"))
		return
	}

	resp, err := resources.NewFromTaskInfo(inferInfo)
	if err != nil {
		log.Errorc(ctx, "new task info resource error: %v", err)
		response.WriteError(ctx.Writer, exception.New(500, "Internal Server Error"))
		return
	}

	response.Write(ctx.Writer, response.Ok(resp))
}

// FilterInference
//
//	@ID				filter-inference
//	@Summary 		Filter inferences
//	@Description
//	@Tags			InferenceController
//	@Accept			json
//	@Produce		json
//	@Param			filter	body    	requests.FilterInferenceRequest     true        "Task IDs"
//	@Success		200		{object}	response.Response{data=[]resources.Inference}
//	@Success		400		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/infer [get]
func (c *InferenceController) FilterInference(ctx *gin.Context) {

	var req requests.FilterInferenceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.WriteError(ctx.Writer, exception.New(40000, err.Error()))
		return
	}

	infoList, err := c.inferenceService.FilterInference(ctx, req.IDs)
	if err != nil {
		log.Errorc(ctx, "new task info resource error: %v", err)
		response.WriteError(ctx.Writer, exception.New(500, "Internal Server Error"))
		return
	}

	resp, err := resources.NewFromTaskInfoList(infoList)
	if err != nil {
		log.Errorc(ctx, "new task info resource error: %v", err)
		response.WriteError(ctx.Writer, exception.New(500, "Internal Server Error"))
		return
	}

	response.Write(ctx.Writer, response.Ok(resp))
}

// CreateInference
//
//	@ID				create-inference
//	@Summary 		Create an inference task
//	@Description
//	@Tags			InferenceController
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			model					formData		string			true	"AI model ID" default(amnd_general)
//	@Param			type					formData		string			false	"Infer type" Enums(t2v,i2v,v2v,upscale) default(t2v)
//	@Param			prompt					formData		string			true	"The prompt or prompts to guide image generation." minlength(1) maxlength(1000)
//	@Param			negative_prompt			formData		string			false	"The prompt or prompts to guide what to not include in image generation."
//	@Param			num_inference_steps		formData		int				false	"More steps usually lead to a higher quality image at the expense of slower inference" default(4) minimum(1) maximum(200)
//	@Param			num_frames				formData		int				false	"The number of video frames to generate. Default FPS: 8" default(16) minimum(16) maximum(32)
//	@Param			width					formData		int				false	"The width in pixels of the generated image/video." default(512)
//	@Param			height					formData		int				false	"The height in pixels of the generated image/video." default(512)
//	@Param			guidance_scale			formData		float32			false	"A higher guidance scale value encourages the model to generate images closely linked to the `prompt` at the expense of lower image quality." default(1.5) minimum(0) maximum(100)
//	@Success		201		{object}	response.Response{data=resources.Inference}
//	@Failure		400		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/infer [post]
func (c *InferenceController) CreateInference(ctx *gin.Context) {

	var req requests.CreateInferenceRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.WriteError(ctx.Writer, exception.New(40000, err.Error()))
		return
	}

	modelProps, exist := c.modelProps.DataMap[req.Model]
	if !exist {
		response.WriteError(ctx.Writer, exception.New(40001, "Model not supported"))
		return
	}
	req.Model = modelProps.Path
	if len(modelProps.TriggerWords) > 1 {
		req.Prompt = fmt.Sprintf("%s %s", req.Prompt, modelProps.TriggerWords)
	}

	resp, err := c.inferenceService.CreateInference(ctx, req)
	if err != nil {
		log.Errorc(ctx, "%v", err)
		response.WriteError(ctx.Writer, exception.New(50000, "Internal Server Error"))
		return
	}

	response.Write(ctx.Writer, response.Created(resp))
}

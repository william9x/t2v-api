package controllers

import (
	"github.com/Braly-Ltd/t2v-api-public/requests"
	goaway "github.com/TwiN/go-away"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/exception"
	"github.com/golibs-starter/golib/web/response"
)

type NotificationController struct {
	profanityDetector *goaway.ProfanityDetector
}

func NewNotificationController(
	profanityDetector *goaway.ProfanityDetector,
) *NotificationController {
	return &NotificationController{
		profanityDetector: profanityDetector,
	}
}

// Subscribe
//
//	@ID				subscribe
//	@Summary 		Subscribe for notifications
//	@Description
//	@Tags			NotificationController
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			user_id				formData		string			true	"User's device ID or user's account (not anonymous) ID"
//	@Param			user_token			formData		string			true	"User's registration token for push notifications"
//	@Param			token_provider		formData		string			true	"Notification provider" default(firebase)
//	@Success		201		{object}	response.Response{data=requests.CreateSubscriptionRequest}
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/noti/subs [post]
func (c *NotificationController) Subscribe(ctx *gin.Context) {
	var req requests.CreateSubscriptionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.WriteError(ctx.Writer, exception.New(40000, err.Error()))
		return
	}

	response.Write(ctx.Writer, response.Created(req))
}

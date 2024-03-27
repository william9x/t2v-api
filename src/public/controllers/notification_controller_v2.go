package controllers

import (
	"github.com/Braly-Ltd/t2v-api-core/ports"
	"github.com/Braly-Ltd/t2v-api-public/requests"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/exception"
	"github.com/golibs-starter/golib/web/response"
)

type NotificationControllerV2 struct {
	notiSubscriptionPort ports.NotificationSubscriptionPort
}

func NewNotificationControllerV2(
	notiSubscriptionPort ports.NotificationSubscriptionPort,
) *NotificationControllerV2 {
	return &NotificationControllerV2{
		notiSubscriptionPort: notiSubscriptionPort,
	}
}

// Subscribe
//
//	@ID				subscribe-v2
//	@Summary 		Subscribe for notifications
//	@Description
//	@Tags			NotificationControllerV2
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			user_id				formData		string			true	"User's device ID or user's account (not anonymous) ID"
//	@Param			user_token			formData		string			true	"User's registration token for push notifications"
//	@Param			token_provider		formData		string			true	"Notification provider" default(firebase)
//	@Success		201		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/api/v2/noti/subs [post]
func (c *NotificationControllerV2) Subscribe(ctx *gin.Context) {
	var req requests.CreateSubscriptionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.WriteError(ctx.Writer, exception.New(40000, err.Error()))
		return
	}

	exist, err := c.notiSubscriptionPort.TokenExist(ctx, req.UserToken)
	if err != nil {
		response.WriteError(ctx.Writer, exception.New(50000, "Failed to check token existence."))
		return
	}

	if exist {
		response.WriteError(ctx.Writer, exception.New(40000, "Token already exists."))
		return
	}

	if _, err := c.notiSubscriptionPort.Subscribe(ctx, req.UserID, req.UserToken, req.TokenProvider); err != nil {
		response.WriteError(ctx.Writer, exception.New(50000, "Failed to subscribe for notifications."))
		return
	}

	response.Write(ctx.Writer, response.Created(nil))
}

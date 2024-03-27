package firebase

import (
	"context"
	"firebase.google.com/go/messaging"
	"fmt"
)

type MessagingClient struct {
	Android   *messaging.Client
	AndroidV2 *messaging.Client
	IOS       *messaging.Client
	IOSV2     *messaging.Client
}

func NewFirebaseMessagingClient(app *Application) (*MessagingClient, error) {

	androidMsgClient, err := app.Android.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase msg client for android: %v", err)
	}

	androidMsgClientV2, err := app.AndroidV2.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase msg client for android: %v", err)
	}

	iosMsgClient, err := app.IOS.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase msg client for ios: %v", err)
	}

	iosV2MsgClient, err := app.IOSV2.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase msg client for ios: %v", err)
	}

	return &MessagingClient{
		Android:   androidMsgClient,
		AndroidV2: androidMsgClientV2,
		IOS:       iosMsgClient,
		IOSV2:     iosV2MsgClient,
	}, nil
}

func (r *MessagingClient) SendNoti(ctx context.Context, agent, taskID, title, body, image, token string) (string, error) {
	var msgClient *messaging.Client
	if agent == "ios" {
		msgClient = r.IOSV2
	} else if agent == "android" {
		msgClient = r.AndroidV2
	}

	data := make(map[string]string)
	data["task_id"] = taskID
	data["title"] = title
	data["body"] = body
	data["image"] = image
	msg := messaging.Message{
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: image,
		},
		Data:  data,
		Token: token,
	}
	return msgClient.Send(ctx, &msg)
}

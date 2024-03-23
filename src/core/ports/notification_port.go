package ports

import (
	"context"
)

type NotificationPort interface {
	SendNoti(ctx context.Context, agent, taskID, title, body, image, token string) (string, error)
}

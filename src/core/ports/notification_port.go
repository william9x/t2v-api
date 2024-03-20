package ports

import (
	"context"
)

type NotificationPort interface {
	SendNoti(ctx context.Context, agent, title, body, image, token string) (string, error)
}

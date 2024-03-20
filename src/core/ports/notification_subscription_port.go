package ports

import (
	"context"
)

type NotificationSubscriptionPort interface {
	Subscribe(ctx context.Context, userID, token, provider string) (interface{}, error)
	TokenExist(ctx context.Context, token string) (bool, error)
}

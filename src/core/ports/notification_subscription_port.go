package ports

import (
	"context"
	"github.com/Braly-Ltd/t2v-api-core/entities"
)

type NotificationSubscriptionPort interface {
	Subscribe(ctx context.Context, userID, token, provider string) (interface{}, error)
	TokenExist(ctx context.Context, token string) (bool, error)
	FindByUserID(ctx context.Context, userID string) ([]entities.NotiSubscription, error)
}

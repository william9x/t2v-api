package adapter

import (
	"context"
	"fmt"
	"github.com/Braly-Ltd/t2v-api-adapter/properties"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NotificationSubscriptionAdapter ...
type NotificationSubscriptionAdapter struct {
	props  *properties.MongoProperties
	client *mongo.Client
}

const (
	databaseName        = "sira"
	notiSubscriptions   = "noti_subscriptions"
	inferenceCollection = "inferences"
)

// NewNotificationSubscriptionAdapter ...
func NewNotificationSubscriptionAdapter(props *properties.MongoProperties, client *mongo.Client) *NotificationSubscriptionAdapter {
	return &NotificationSubscriptionAdapter{props: props, client: client}
}

func (r *NotificationSubscriptionAdapter) Subscribe(ctx context.Context, userID, token, provider string) (interface{}, error) {
	collection := r.client.Database(databaseName).Collection(notiSubscriptions)

	doc := bson.M{"user_id": userID, "token": token, "provider": provider}
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("error while inserting subscription: %w", err)
	}

	return result.InsertedID, nil
}

func (r *NotificationSubscriptionAdapter) TokenExist(ctx context.Context, token string) (bool, error) {
	collection := r.client.Database(databaseName).Collection(notiSubscriptions)

	countOpts := new(options.CountOptions)
	countOpts.SetLimit(1)
	count, err := collection.CountDocuments(ctx, bson.M{"token": token}, countOpts)
	if err != nil {
		return false, fmt.Errorf("error while counting subcriptions: %w", err)
	}
	return count > 0, nil
}

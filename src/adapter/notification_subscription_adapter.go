package adapter

import (
	"context"
	"fmt"
	"github.com/Braly-Ltd/t2v-api-adapter/properties"
	"github.com/Braly-Ltd/t2v-api-core/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NotificationSubscriptionAdapter ...
type NotificationSubscriptionAdapter struct {
	props  *properties.MongoProperties
	client *mongo.Client
}

// NewNotificationSubscriptionAdapter ...
func NewNotificationSubscriptionAdapter(props *properties.MongoProperties, client *mongo.Client) *NotificationSubscriptionAdapter {
	return &NotificationSubscriptionAdapter{props: props, client: client}
}

func (r *NotificationSubscriptionAdapter) Subscribe(ctx context.Context, userID, token, provider string) (interface{}, error) {
	doc := bson.M{"user_id": userID, "token": token, "provider": provider}
	result, err := r.collection().InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("error while inserting subscription: %w", err)
	}

	return result.InsertedID, nil
}

func (r *NotificationSubscriptionAdapter) TokenExist(ctx context.Context, token string) (bool, error) {
	countOpts := new(options.CountOptions)
	countOpts.SetLimit(1)
	count, err := r.collection().CountDocuments(ctx, bson.M{"token": token}, countOpts)
	if err != nil {
		return false, fmt.Errorf("error while counting subcriptions: %w", err)
	}
	return count > 0, nil
}

func (r *NotificationSubscriptionAdapter) FindByUserID(ctx context.Context, userID string) ([]entities.NotiSubscription, error) {
	result, err := r.collection().Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("error while fetching user_id: %w", err)
	}

	var subscriptions []entities.NotiSubscription
	if err := result.All(ctx, &subscriptions); err != nil {
		return nil, fmt.Errorf("error while decoding user_id: %w", err)
	}

	return subscriptions, nil
}

func (r *NotificationSubscriptionAdapter) collection() *mongo.Collection {
	return r.client.Database(mongodbDatabaseName).Collection(mongodbNotiSubCollection)
}

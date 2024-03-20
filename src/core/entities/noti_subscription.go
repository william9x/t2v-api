package entities

type NotiSubscription struct {
	UserID        string `bson:"user_id"`
	UserToken     string `bson:"token"`
	TokenProvider string `bson:"provider"`
}

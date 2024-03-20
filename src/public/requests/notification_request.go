package requests

// CreateSubscriptionRequest ...
type CreateSubscriptionRequest struct {
	UserID        string `json:"user_id" form:"user_id" binding:"notblank"`
	UserToken     string `json:"user_token" form:"user_token" binding:"notblank"`
	TokenProvider string `json:"token_provider" form:"token_provider,default=firebase" binding:"notblank"`
}

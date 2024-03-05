package adapter

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/Braly-Ltd/t2v-api-core/entities"
	"github.com/golibs-starter/golib/log"
)

// FirebaseAdapter ...
type FirebaseAdapter struct {
	authClient *auth.Client
}

// NewFirebaseAdapter ...
func NewFirebaseAdapter(authClient *auth.Client) (*FirebaseAdapter, error) {
	return &FirebaseAdapter{
		authClient: authClient,
	}, nil
}

func (r *FirebaseAdapter) Authenticate(ctx context.Context, token string) (entities.TokenData, error) {
	tokenData, err := r.authClient.VerifyIDToken(ctx, token)
	if err != nil {
		log.Warnf("error verifying token: %v", err)
		return entities.TokenData{}, err
	}
	return entities.TokenData{
		Issuer:   tokenData.Issuer,
		Expires:  tokenData.Expires,
		IssuedAt: tokenData.IssuedAt,
		Subject:  tokenData.Subject,
		UserID:   tokenData.UID,
		Claims:   tokenData.Claims,
	}, nil
}

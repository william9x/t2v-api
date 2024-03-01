package adapter

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/Braly-Ltd/t2v-api-adapter/properties"
	"github.com/Braly-Ltd/t2v-api-core/entities"
	"github.com/golibs-starter/golib/log"
	"google.golang.org/api/option"
)

// FirebaseAdapter ...
type FirebaseAdapter struct {
	authClient *auth.Client
}

// NewFirebaseAdapter ...
func NewFirebaseAdapter(props *properties.FirebaseProperties) (*FirebaseAdapter, error) {
	opt := option.WithCredentialsFile(props.CredentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase auth client: %v", err)
	}
	return &FirebaseAdapter{
		authClient: client,
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

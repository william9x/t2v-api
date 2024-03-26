package firebase

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/Braly-Ltd/t2v-api-core/entities"
	"github.com/golibs-starter/golib/log"
)

type AuthClient struct {
	Android   *auth.Client
	AndroidV2 *auth.Client
	IOS       *auth.Client
	IOSV2     *auth.Client
}

func NewFirebaseAuthClient(app *Application) (*AuthClient, error) {

	androidAuthClient, err := app.Android.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase auth client for android: %v", err)
	}

	androidAuthClientV2, err := app.AndroidV2.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase auth client for android: %v", err)
	}

	iosAuthClient, err := app.IOS.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase auth client for ios: %v", err)
	}

	iosAuthClientV2, err := app.IOS.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase auth client for ios V2: %v", err)
	}

	return &AuthClient{
		Android:   androidAuthClient,
		AndroidV2: androidAuthClientV2,
		IOS:       iosAuthClient,
		IOSV2:     iosAuthClientV2,
	}, nil
}

func (r *AuthClient) Authenticate(ctx context.Context, agent, token string) (entities.TokenData, error) {

	var tokenData *auth.Token
	var err error

	if agent == "ios" {
		tokenData, err = r.IOS.VerifyIDToken(ctx, token)
	} else if agent == "iosV2" {
		tokenData, err = r.IOSV2.VerifyIDToken(ctx, token)
	} else if agent == "android" {
		tokenData, err = r.Android.VerifyIDToken(ctx, token)
	} else if agent == "androidV2" {
		tokenData, err = r.AndroidV2.VerifyIDToken(ctx, token)
	} else {
		log.Warnf("invalid agent, request origin could not be from app: %s", agent)
		return entities.TokenData{}, fmt.Errorf("invalid agent: %s", agent)
	}

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

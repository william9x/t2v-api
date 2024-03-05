package clients

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/Braly-Ltd/t2v-api-adapter/properties"
	"google.golang.org/api/option"
)

func NewFirebaseAuthClient(props *properties.FirebaseProperties) (*auth.Client, error) {
	opt := option.WithCredentialsFile(props.CredentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}
	return app.Auth(context.Background())
}

package firebase

import (
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/Braly-Ltd/t2v-api-adapter/properties"
	"google.golang.org/api/option"
)

type Application struct {
	Android   *firebase.App
	AndroidV2 *firebase.App
	IOS       *firebase.App
	IOSV2     *firebase.App
}

func NewFirebaseApplication(props *properties.FirebaseProperties) (*Application, error) {
	androidApp, err := newFirebaseApp(props.CredentialsFileAndroid)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app for android: %v", err)
	}

	androidAppV2, err := newFirebaseApp(props.CredentialsFileAndroidV2)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app for android V2: %v", err)
	}

	iosApp, err := newFirebaseApp(props.CredentialsFileIOS)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app for ios: %v", err)
	}

	iosAppV2, err := newFirebaseApp(props.CredentialsFileIOSV2)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app for ios V2: %v", err)
	}

	return &Application{
		Android:   androidApp,
		AndroidV2: androidAppV2,
		IOS:       iosApp,
		IOSV2:     iosAppV2,
	}, nil
}

func newFirebaseApp(credentialsFile string) (*firebase.App, error) {
	opt := option.WithCredentialsFile(credentialsFile)
	return firebase.NewApp(context.Background(), nil, opt)
}

package firebase

import (
	"fmt"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type Firebase struct {
	app *firebase.App
}

func New(firebaseCredentialsPath string) (*Firebase, error) {
	opt := option.WithCredentialsFile(firebaseCredentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return &Firebase{
		app: app,
	}, nil
}

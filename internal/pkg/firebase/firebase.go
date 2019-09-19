package firebase

import (
	"golang.org/x/net/context"

	goFirebase "firebase.google.com/go"
	goFirebaseAuth "firebase.google.com/go/auth"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

type Firebase struct {
	app    *goFirebase.App
	client *goFirebaseAuth.Client
}

func New(firebaseCredentialsPath string) (*Firebase, error) {
	opt := option.WithCredentialsFile(firebaseCredentialsPath)
	app, err := goFirebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Error().Err(err).Msg("initialising firebase app")
		return nil, ErrUnexpected{}
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("getting firebase auth client")
		return nil, ErrUnexpected{}
	}

	return &Firebase{
		app:    app,
		client: client,
	}, nil
}

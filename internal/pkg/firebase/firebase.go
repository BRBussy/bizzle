package firebase

import (
	"golang.org/x/net/context"
	"time"

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
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	app, err := goFirebase.NewApp(ctx, nil, opt)
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

func (f *Firebase) GetUserByEmail(email string) (*goFirebaseAuth.UserRecord, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	user, err := f.client.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Msg("getting user by email")
		return nil, ErrUnexpected{}
	}
	return user, nil
}

func (f *Firebase) CreateUser(firebaseUserDetails *goFirebaseAuth.UserToCreate) (*goFirebaseAuth.UserRecord, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	user, err := f.client.CreateUser(ctx, firebaseUserDetails)
	if err != nil {
		log.Error().Err(err).Msg("creating user")
		return nil, ErrUnexpected{}
	}
	return user, nil
}

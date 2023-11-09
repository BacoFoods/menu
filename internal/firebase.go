package internal

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type Firebase struct {
	*db.Client
	app *firebase.App
}

func MustNewFirebase(auth, dbURL string) *Firebase {
	fs, err := NewFirebase(auth, dbURL)
	if err != nil {
		panic(err)
	}

	return fs
}

func NewFirebase(auth, dbURL string) (*Firebase, error) {
	ctx := context.Background()
	creds, err := shared.Base64Decode(auth)
	if err != nil {
		return nil, err
	}

	logrus.Info("Connecting to firebase:", dbURL)

	opt := option.WithCredentialsJSON(creds)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	db, err := app.DatabaseWithURL(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	return &Firebase{Client: db, app: app}, nil
}

package internal

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/BacoFoods/menu/pkg/shared"
	"google.golang.org/api/option"
)

type Firebase struct {
	app *firebase.App
}

func NewFirebase() (*Firebase, error) {
	ctx := context.Background()
	creds, err := shared.Base64Decode(Config.FirestoreConfig.AuthB64)
	if err != nil {
		return nil, err
	}

	opt := option.WithCredentialsJSON(creds)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	return &Firebase{app}, nil
}

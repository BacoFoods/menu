package internal

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type Firestore struct {
	Config FirestoreConfig
	Client *firestore.Client
}

func NewFirestore(ifc FirestoreConfig) (*Firestore, error) {
	ctx := context.Background()
	auth, err := shared.Base64Decode(ifc.AuthB64)
	if err != nil {
		logrus.Errorf("Failed to Marshal firestore credentials: %v", err)
		return nil, err
	}

	credOpt := option.WithCredentialsJSON(auth)
	client, err := firestore.NewClient(ctx, ifc.ProjectID, credOpt)
	if err != nil {
		logrus.Errorf("Failed to create Firestore client: %v", err)
		return nil, err
	}

	return &Firestore{
		Config: ifc,
		Client: client,
	}, nil
}

package internal

import (
	"context"
	"encoding/base64"

	"cloud.google.com/go/firestore"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type Firestore struct {
	Config FirestoreConfig
	Client *firestore.Client
}

func NewFirestore(ifc FirestoreConfig) (*Firestore, error) {
	ctx := context.Background()

	// Decode base64 auth
	decoded := base64.StdEncoding.DecodedLen(len(ifc.AuthB64))
	decodedString := make([]byte, decoded)
	n, err := base64.StdEncoding.Decode(decodedString, []byte(ifc.AuthB64))
	if err != nil {
		logrus.Errorf("Failed to Marshal firestore credentials: %v", err)
		return nil, err
	}

	credOpt := option.WithCredentialsJSON(decodedString[:n])
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

package firebase

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/construction_crm/internal/construction_crm/types/config"
	"google.golang.org/api/option"
)

type FireBase struct {
	client *messaging.Client
}

func NewFireBase(cfg *config.Fire) (*FireBase, error) {
	sa := option.WithCredentialsFile(cfg.Path)
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}

	return &FireBase{
		client: client,
	}, nil
}

func (fb *FireBase) Send(message *messaging.Message) error {

	_, err := fb.client.Send(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}

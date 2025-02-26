package auth

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"github.com/Financial-Partner/server/internal/config"
)

//go:generate mockgen -source=firebase.go -destination=mocks/firebase_mock.go -package=mocks

type FirebaseAuth interface {
	VerifyIDToken(ctx context.Context, idToken string) (*Token, error)
}

type Token = auth.Token

type Client struct {
	auth FirebaseAuth
}

func NewWithAuth(auth FirebaseAuth) *Client {
	return &Client{auth: auth}
}

func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	opt := option.WithCredentialsFile(cfg.Firebase.CredentialFile)
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: cfg.Firebase.ProjectID,
	}, opt)
	if err != nil {
		return nil, err
	}

	firebaseAuth, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return NewWithAuth(firebaseAuth), nil
}

func (c *Client) VerifyToken(ctx context.Context, idToken string) (*Token, error) {
	return c.auth.VerifyIDToken(ctx, idToken)
}

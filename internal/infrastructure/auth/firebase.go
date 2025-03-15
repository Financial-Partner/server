package auth

import (
	"context"
	"errors"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"github.com/Financial-Partner/server/internal/config"
)

var (
	// ErrBypassTokenDisabled is returned when bypass token is used but the feature is disabled
	ErrBypassTokenDisabled = errors.New("bypass token feature is disabled")
)

//go:generate mockgen -source=firebase.go -destination=firebase_mock.go -package=auth

type FirebaseAuth interface {
	VerifyIDToken(ctx context.Context, idToken string) (*Token, error)
}

type Token = auth.Token

type Client struct {
	auth          FirebaseAuth
	bypassToken   string
	bypassEnabled bool
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

	return &Client{
		auth:          firebaseAuth,
		bypassToken:   cfg.Firebase.BypassToken,
		bypassEnabled: cfg.Firebase.BypassEnabled,
	}, nil
}

func (c *Client) VerifyToken(ctx context.Context, idToken string) (*Token, error) {
	if c.bypassEnabled && c.bypassToken != "" && idToken == c.bypassToken {
		return &Token{
			UID: "bypass-user-id",
			Claims: map[string]interface{}{
				"email": "bypass@example.com",
				"name":  "Bypass User",
			},
		}, nil
	}

	// If it's not the bypass token, proceed with normal verification
	return c.auth.VerifyIDToken(ctx, idToken)
}

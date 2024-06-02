package firebase

import (
	"context"
	"os"

	"Alice-Seahat-Healthcare/seahat-be/apperror"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type Firebase interface {
	GetAuthIdentity(ctx context.Context, googleToken string) (map[string]string, error)
}

type firebaseImpl struct {
	AuthClient *auth.Client
}

func New() (*firebaseImpl, error) {
	wd, _ := os.Getwd()
	credentialPath := wd + "/firebase.json"
	ctx := context.Background()

	opt := option.WithCredentialsFile(credentialPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &firebaseImpl{
		AuthClient: authClient,
	}, nil
}

func (f *firebaseImpl) GetAuthIdentity(ctx context.Context, googleToken string) (map[string]string, error) {
	t, err := f.AuthClient.VerifyIDToken(ctx, googleToken)
	if err != nil {
		return nil, apperror.ErrInvalidToken
	}

	return map[string]string{
		"Name":    t.Claims["name"].(string),
		"Email":   t.Claims["email"].(string),
		"Picture": t.Claims["picture"].(string),
	}, nil
}

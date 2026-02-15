package authenticator

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"rminder/internal/pkg/config"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config

	cfg     config.AuthConfig
	once    sync.Once
	initErr error
}

// New instantiates the *Authenticator without connecting to the OIDC provider.
// The provider is lazily initialized on first use.
func New(cfg config.AuthConfig) *Authenticator {
	return &Authenticator{
		cfg: cfg,
	}
}

// Init ensures the OIDC provider is initialized. Safe for concurrent use.
func (a *Authenticator) Init() error {
	a.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		provider, err := oidc.NewProvider(
			ctx,
			"https://"+a.cfg.Domain+"/",
		)
		if err != nil {
			a.initErr = err
			return
		}

		a.Provider = provider
		a.Config = oauth2.Config{
			ClientID:     a.cfg.ClientID,
			ClientSecret: a.cfg.ClientSecret,
			RedirectURL:  a.cfg.CallbackUrl,
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile"},
		}
	})

	return a.initErr
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

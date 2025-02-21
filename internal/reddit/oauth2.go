package reddit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cenkalti/backoff/v5"
	"golang.org/x/oauth2"
)

// NewOAuth2Client creates and configures an OAuth2 HTTP client with retry logic
// for rate limiting.
func NewOAuth2Client(ctx context.Context, cfg Config) (*http.Client, error) {
	// Create a custom HTTP client with the User-Agent header, used for OAuth2
	// token requests.
	// Ref: https://github.com/golang/oauth2/issues/179
	client := &http.Client{
		Transport: &userAgentTransport{
			base:      http.DefaultTransport,
			userAgent: cfg.UserAgent,
		},
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, client)

	oauthCfg := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL:  fmt.Sprintf("%s/api/v1/access_token", cfg.BaseURL),
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}

	// Get the initial token.
	tok, err := getToken(ctx, oauthCfg, cfg.Username, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("error getting initial token: %w", err)
	}

	// Create a new OAuth2 client with the initial token.
	return oauthCfg.Client(ctx, tok), nil
}

// getToken attempts to get an OAuth2 token with retry logic for rate limiting.
func getToken(ctx context.Context, oauthCfg *oauth2.Config, username, password string) (*oauth2.Token, error) {
	operation := func() (*oauth2.Token, error) {
		tok, err := oauthCfg.PasswordCredentialsToken(ctx, username, password)
		if err != nil {
			var retrieveError *oauth2.RetrieveError
			if errors.As(err, &retrieveError) {
				// If we are being rate limited, use the Retry-After header.
				if retrieveError.Response.StatusCode == 429 {
					retryAfter := retrieveError.Response.Header.Get("Retry-After")
					if seconds, err := strconv.Atoi(retryAfter); err == nil {
						return nil, backoff.RetryAfter(seconds)
					}
					// If no valid Retry-After header, just return the error to
					// use default backoff.
					return nil, err
				}
			}
			// For other errors, stop retrying.
			return nil, backoff.Permanent(err)
		}
		return tok, nil
	}

	// Request a token. Retry up to 5 times with exponential backoff if we are
	// being rate limited.
	// TODO: make the max tries configurable -2025-02-21
	tok, err := backoff.Retry(
		ctx,
		operation,
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxTries(5),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting token after retries: %w", err)
	}
	return tok, nil
}

// userAgentTransport adds a User-Agent header to all requests.
type userAgentTransport struct {
	base      http.RoundTripper
	userAgent string
}

func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.userAgent)
	return t.base.RoundTrip(req)
}

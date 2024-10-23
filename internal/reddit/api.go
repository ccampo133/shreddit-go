package reddit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
)

const (
	DefaultBaseURL   = "https://oauth.reddit.com"
	DefaultUserAgent = "shreddit-go"
)

// TODO: doc -ccampo 2024-10-22
type Params struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	UserAgent    string
}

// TODO: doc -ccampo 2024-10-22
type Client struct {
	rc *resty.Client
}

// TODO: doc -ccampo 2024-10-22
func NewClient(ctx context.Context, params Params) (*Client, error) {
	// Set defaults.
	if params.BaseURL == "" {
		params.BaseURL = DefaultBaseURL
	}
	if params.UserAgent == "" {
		params.UserAgent = DefaultUserAgent
	}

	// Configure auto-refreshing OAuth2 client.
	conf := &oauth2.Config{
		ClientID:     params.ClientID,
		ClientSecret: params.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: fmt.Sprintf("%s/api/v1/access_token", params.BaseURL),
		},
	}
	tok, err := conf.PasswordCredentialsToken(ctx, params.Username, params.Password)
	if err != nil {
		return nil, fmt.Errorf("error getting token: %w", err)
	}

	rc := resty.NewWithClient(conf.Client(ctx, tok)).
		SetBaseURL(params.BaseURL).
		SetHeader("User-Agent", params.UserAgent)
	return &Client{rc: rc}, nil
}

// TODO: doc -ccampo 2024-10-22
func (c *Client) GetComments(username string) ([]Comment, error) {
	resp, err := c.rc.R().Get(fmt.Sprintf("/user/%s/comments.json", username))
	if err != nil {
		return nil, fmt.Errorf("error getting comments: %w", err)
	}
	var body Listing[Comment]
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	// TODO: pagination -ccampo 2024-10-22
	return body.Items(), nil
}

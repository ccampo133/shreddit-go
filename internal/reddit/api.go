package reddit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
)

const (
	defaultBaseURL   = "https://oauth.reddit.com"
	defaultUserAgent = "shreddit-go"

	// Reddit "things" (e.g. comments, posts) are identified by a prefix followed
	// by an ID. See https://www.reddit.com/dev/api/#fullnames.
	commentPrefix = "t1_"
	postPrefix    = "t3_"
)

var (
	ErrRateLimited = fmt.Errorf("rate limited")
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
		params.BaseURL = defaultBaseURL
	}
	if params.UserAgent == "" {
		params.UserAgent = defaultUserAgent
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
		return nil, fmt.Errorf("error unmarshalling comments response: %w", err)
	}
	// TODO: pagination -ccampo 2024-10-22
	return body.Items(), nil
}

// TODO: doc -ccampo 2024-10-25
func (c *Client) EditComment(id, body string) error {
	fullName := commentPrefix + id
	resp, err := c.rc.R().
		SetQueryParams(map[string]string{"raw_json": "1"}).
		SetFormData(map[string]string{"thing_id": fullName, "text": body}).
		Post("/api/editusertext")
	if err != nil {
		return fmt.Errorf("error editing comment with id %s: %w", fullName, err)
	}
	var editResp EditResponse
	if err := json.Unmarshal(resp.Body(), &editResp); err != nil {
		return fmt.Errorf("error unmarshalling edit response: %w", err)
	}
	if !editResp.Success {
		if editResp.IsRateLimited() {
			return ErrRateLimited
		}
		// TODO: log the response -ccampo 2024-10-25
		return fmt.Errorf("API failure editing comment with id %s", fullName)
	}
	return nil
}

// TODO: doc -ccampo 2024-10-25
func (c *Client) DeleteComment(id string) error {
	fullName := commentPrefix + id
	if err := c.deleteThing(fullName); err != nil {
		return fmt.Errorf("error deleting comment with id %s: %w", fullName, err)
	}
	return nil
}

// TODO: doc -ccampo 2024-10-25
func (c *Client) DeletePost(id string) error {
	fullName := postPrefix + id
	if err := c.deleteThing(fullName); err != nil {
		return fmt.Errorf("error deleting post with id %s: %w", fullName, err)
	}
	return nil
}

// TODO: doc -ccampo 2024-10-25
// See https://www.reddit.com/dev/api/#fullnames regarding the fullname format.
func (c *Client) deleteThing(fullName string) error {
	_, err := c.rc.R().
		SetFormData(map[string]string{"id": fullName}).
		Post("/api/del")
	return err
}

package reddit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	defaultBaseURL   = "https://oauth.reddit.com"
	defaultUserAgent = "shreddit-go"
)

var (
	ErrRateLimited = fmt.Errorf("rate limited")
)

// TODO: doc -2024-10-22
type Config struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	UserAgent    string
}

// TODO: doc -2024-10-22
type Client struct {
	rc *resty.Client
}

// TODO: doc -2024-10-22
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	// Set defaults.
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = defaultUserAgent
	}

	httpClient, err := NewOAuth2Client(ctx, cfg)
	if err != nil {
		return nil, err
	}

	rc := resty.NewWithClient(httpClient).
		SetBaseURL(cfg.BaseURL).
		SetHeader("User-Agent", cfg.UserAgent)
	return &Client{rc: rc}, nil
}

// TODO: doc -2024-10-30
func (c *Client) GetPosts(username, after string) (*Listing[Post], error) {
	req := c.rc.R()
	if after != "" {
		req.SetQueryParam("after", after)
	}
	resp, err := req.Get(fmt.Sprintf("/user/%s/submitted.json", username))
	if err != nil {
		return nil, fmt.Errorf("error getting posts: %w", err)
	}
	var body Listing[Post]
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling post listing: %w", err)
	}
	return &body, nil
}

// TODO: doc -2024-10-30
func (c *Client) GetSavedPosts(username, after string) (*Listing[Comment], error) {
	req := c.rc.R().
		SetQueryParams(map[string]string{"type": "links"})
	if after != "" {
		req.SetQueryParam("after", after)
	}
	resp, err := req.Get(fmt.Sprintf("/user/%s/saved.json", username))
	if err != nil {
		return nil, fmt.Errorf("error getting saved posts: %w", err)
	}
	var body Listing[Comment]
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling post listing: %w", err)
	}
	return &body, nil
}

// TODO: doc -2024-10-22
func (c *Client) GetComments(username, after string) (*Listing[Comment], error) {
	req := c.rc.R()
	if after != "" {
		req.SetQueryParam("after", after)
	}
	resp, err := req.Get(fmt.Sprintf("/user/%s/comments.json", username))
	if err != nil {
		return nil, fmt.Errorf("error getting comments: %w", err)
	}
	var body Listing[Comment]
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling comment listing: %w", err)
	}
	return &body, nil
}

// TODO: doc -2024-10-30
func (c *Client) GetSavedComments(username, after string) (*Listing[Comment], error) {
	req := c.rc.R().
		SetQueryParams(map[string]string{"type": "comments"})
	if after != "" {
		req.SetQueryParam("after", after)
	}
	resp, err := req.Get(fmt.Sprintf("/user/%s/saved.json", username))
	if err != nil {
		return nil, fmt.Errorf("error getting saved comments: %w", err)
	}
	var body Listing[Comment]
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling comment listing: %w", err)
	}
	return &body, nil
}

// TODO: doc -2024-10-25
func (c *Client) EditComment(id, body string) error {
	fullName := commentFullName(id)
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
		// TODO: log the response -2024-10-25
		return fmt.Errorf("API failure editing comment with id %s", fullName)
	}
	return nil
}

// TODO: doc -2024-10-30
func (c *Client) UnsaveComment(id string) error {
	fullName := commentFullName(id)
	if err := c.unsaveThing(fullName); err != nil {
		return fmt.Errorf("error unsaving comment with id %s: %w", fullName, err)
	}
	return nil
}

// TODO: doc -2024-10-30
func (c *Client) UnsavePost(id string) error {
	fullName := postFullName(id)
	if err := c.unsaveThing(fullName); err != nil {
		return fmt.Errorf("error unsaving post with id %s: %w", fullName, err)
	}
	return nil
}

// TODO: doc -2024-10-25
func (c *Client) DeleteComment(id string) error {
	fullName := commentFullName(id)
	if err := c.deleteThing(fullName); err != nil {
		return fmt.Errorf("error deleting comment with id %s: %w", fullName, err)
	}
	return nil
}

// TODO: doc -2024-10-25
func (c *Client) DeletePost(id string) error {
	fullName := postFullName(id)
	if err := c.deleteThing(fullName); err != nil {
		return fmt.Errorf("error deleting post with id %s: %w", fullName, err)
	}
	return nil
}

// TODO: doc -2024-10-30
func (c *Client) unsaveThing(fullName string) error {
	_, err := c.rc.R().
		SetFormData(map[string]string{"id": fullName}).
		Post("/api/unsave")
	return err
}

// TODO: doc -2024-10-25
func (c *Client) deleteThing(fullName string) error {
	_, err := c.rc.R().
		SetFormData(map[string]string{"id": fullName}).
		Post("/api/del")
	return err
}

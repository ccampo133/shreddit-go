package reddit

import (
	"encoding/json"
	"fmt"
	"time"
)

const rateLimitErrorText = ".error.RATELIMIT.field-ratelimit"

// TODO: doc -2024-10-22
type Listing[T any] struct {
	Data struct {
		Before   string `json:"before"`
		After    string `json:"after"`
		Children []struct {
			Data T `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// TODO: doc -2024-10-22
func (l *Listing[T]) Items() []T {
	items := make([]T, 0, len(l.Data.Children))
	for _, item := range l.Data.Children {
		items = append(items, item.Data)
	}
	return items
}

// TODO: doc -2024-10-22
type Comment struct {
	ID         string     `json:"id"`
	Body       string     `json:"body"`
	Permalink  string     `json:"permalink"`
	Subreddit  string     `json:"subreddit"`
	Score      int  `json:"score"`
	CreatedUTC Time `json:"created_utc"`
}

// TODO: doc -2024-10-30
type Post struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Permalink  string `json:"permalink"`
	Subreddit  string `json:"subreddit"`
	Score      int    `json:"score"`
	CreatedUTC Time   `json:"created_utc"`
}

// Time is a type used to unmarshal Reddit's weird floating point timestamps.
// Reddit's API returns timestamps as Unix epoch timestamps, but as floating
// point numbers (for some reason). This type is used to unmarshal those
// timestamps into Go's time.Time type.
type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var sec float64
	if err := json.Unmarshal(data, &sec); err != nil {
		return fmt.Errorf("error unmarshalling unix time: %w", err)
	}
	t.Time = time.Unix(int64(sec), 0)
	return nil
}

// EditResponse is the response from the Reddit API when editing a comment. It
// has a weird structure - see the tests for examples.
type EditResponse struct {
	JQuery  []any `json:"jquery"` // ...wat?
	Success bool  `json:"success"`
}

// IsRateLimited checks if the response indicates that the request was rate
// limited.
func (resp *EditResponse) IsRateLimited() bool {
	// Wtf kind of API is this???
	for _, elem := range resp.JQuery {
		if arr, ok := elem.([]any); ok && len(arr) > 3 {
			if word, ok := arr[2].(string); ok && word == "call" {
				if arr2, ok := arr[3].([]any); ok && len(arr2) == 1 {
					if val, ok := arr2[0].(string); ok && val == rateLimitErrorText {
						return true
					}
				}
			}
		}
	}
	return false
}

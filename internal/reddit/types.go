package reddit

import (
	"encoding/json"
	"fmt"
	"time"
)

// TODO: doc -ccampo 2024-10-22
type Listing[T any] struct {
	Data struct {
		Before   string `json:"before"`
		After    string `json:"after"`
		Children []struct {
			Data T `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// TODO: doc -ccampo 2024-10-22
func (l *Listing[T]) Items() []T {
	items := make([]T, 0, len(l.Data.Children))
	for _, item := range l.Data.Children {
		items = append(items, item.Data)
	}
	return items
}

// TODO: doc -ccampo 2024-10-22
type Comment struct {
	ID         string     `json:"id"`
	Body       string     `json:"body"`
	Permalink  string     `json:"permalink"`
	Subreddit  string     `json:"subreddit"`
	Score      int        `json:"score"`
	CreatedUTC redditTime `json:"created_utc"`
}

// Reddit's API returns timestamps as Unix epoch timestamps, but as floating
// point numbers (for some reason). This type is used to unmarshal those
// timestamps into Go's time.Time type.
type redditTime struct {
	time.Time
}

func (t *redditTime) UnmarshalJSON(data []byte) error {
	var sec float64
	if err := json.Unmarshal(data, &sec); err != nil {
		return fmt.Errorf("error unmarshalling unix time: %w", err)
	}
	t.Time = time.Unix(int64(sec), 0)
	return nil
}

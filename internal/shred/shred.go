package shred

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ccampo133/shreddit-go/internal/reddit"
)

const (
	// TODO: doc -2024-10-30
	DefaultReplacementComment = "[deleted]"
)

// TODO: doc -2024-10-30
type Config struct {
	Username           string
	DryRun             bool
	SkipComments       bool
	SkipPosts          bool
	SkipSavedComments  bool
	SkipSavedPosts     bool
	EditOnly           bool
	Before             time.Time
	MaxScore           *int
	MaxDays            *int
	ReplacementComment string
	Sleep              time.Duration
}

// TODO: doc -2024-10-30
type Shredder struct {
	client *reddit.Client
	cfg    Config
}

// TODO: doc -2024-10-30
func NewShredder(client *reddit.Client, cfg Config) *Shredder {
	if cfg.Before.IsZero() {
		cfg.Before = time.Now()
		if cfg.MaxDays != nil {
			cfg.Before = cfg.Before.AddDate(0, 0, -*cfg.MaxDays)
		}
	}
	if cfg.ReplacementComment == "" {
		cfg.ReplacementComment = DefaultReplacementComment
	}
	if cfg.Sleep == 0 {
		cfg.Sleep = 2 * time.Second
	}
	return &Shredder{client: client, cfg: cfg}
}

// TODO: doc -2024-10-30
func (s *Shredder) Shred() error {
	// Comments
	if !s.cfg.SkipComments {
		if err := s.pager(s.shredComments); err != nil {
			return fmt.Errorf("error shredding comments: %w", err)
		}
	}
	// Posts
	if !s.cfg.SkipPosts {
		if err := s.pager(s.shredPosts); err != nil {
			return fmt.Errorf("error shredding posts: %w", err)
		}
	}
	// Saved comments
	if !s.cfg.SkipSavedComments {
		if err := s.pager(s.shredSavedComments); err != nil {
			return fmt.Errorf("error shredding saved comments: %w", err)
		}
	}
	// Saved posts
	if !s.cfg.SkipSavedPosts {
		if err := s.pager(s.shredSavedPosts); err != nil {
			return fmt.Errorf("error shredding saved posts: %w", err)
		}
	}
	return nil
}

// TODO: doc -2024-10-30
func (s *Shredder) shredComments(after string) (string, error) {
	res, err := s.client.GetComments(s.cfg.Username, after)
	if err != nil {
		return "", fmt.Errorf("error getting comments: %w", err)
	}
	for i, comment := range res.Data.Children {
		// Skip comments younger than the cutoff time.
		if comment.Data.CreatedUTC.After(s.cfg.Before) {
			slog.Info(
				"Skipping comment (created after cutoff)",
				"created", comment.Data.CreatedUTC.Time,
				"permalink", comment.Data.Permalink,
			)
			continue
		}
		// Skip comments with a score above the cutoff.
		if s.cfg.MaxScore != nil && comment.Data.Score > *s.cfg.MaxScore {
			slog.Info(
				"Skipping comment (score > max score)",
				"score", comment.Data.Score,
				"permalink", comment.Data.Permalink,
			)
			continue
		}
		// Dry run; just log what we would do.
		if s.cfg.DryRun {
			slog.Info("Would shred comment (dry-run)", "permalink", comment.Data.Permalink)
			continue
		}
		// Edit the comment.
		if err := s.client.EditComment(comment.Data.ID, s.cfg.ReplacementComment); err != nil {
			// TODO: handle rate limiting error -2024-10-31
			return "", fmt.Errorf("error editing comment: %w", err)
		}
		if !s.cfg.EditOnly {
			time.Sleep(s.cfg.Sleep)
			// Delete the comment.
			if err := s.client.DeleteComment(comment.Data.ID); err != nil {
				return "", fmt.Errorf("error deleting comment: %w", err)
			}
		}
		slog.Info("Successfully shredded comment", "permalink", comment.Data.Permalink)
		if i < len(res.Data.Children)-1 {
			time.Sleep(s.cfg.Sleep)
		}
	}
	return res.Data.After, nil
}

func (s *Shredder) shredPosts(after string) (string, error) {
	res, err := s.client.GetPosts(s.cfg.Username, after)
	if err != nil {
		return "", fmt.Errorf("error getting posts: %w", err)
	}
	for i, post := range res.Data.Children {
		// Skip posts younger than the cutoff time.
		if post.Data.CreatedUTC.After(s.cfg.Before) {
			slog.Info(
				"Skipping posts (created after cutoff)",
				"created", post.Data.CreatedUTC.Time,
				"permalink", post.Data.Permalink,
			)
			continue
		}
		// Skip posts with a score above the cutoff.
		if s.cfg.MaxScore != nil && post.Data.Score > *s.cfg.MaxScore {
			slog.Info(
				"Skipping post (score > max score)",
				"score", post.Data.Score,
				"permalink", post.Data.Permalink,
			)
			continue
		}
		// Dry run; just log what we would do.
		if s.cfg.DryRun {
			slog.Info("Would shred post (dry-run)", "permalink", post.Data.Permalink)
			continue
		}
		// Delete the post.
		if err := s.client.DeletePost(post.Data.ID); err != nil {
			return "", fmt.Errorf("error deleting post: %w", err)
		}
		slog.Info("Successfully shredded post", "permalink", post.Data.Permalink)
		if i < len(res.Data.Children)-1 {
			time.Sleep(s.cfg.Sleep)
		}
	}
	return res.Data.After, nil
}

func (s *Shredder) shredSavedComments(after string) (string, error) {
	// TODO: implement -2024-10-30
	return "", nil
}

func (s *Shredder) shredSavedPosts(after string) (string, error) {
	// TODO: implement -2024-10-30
	return "", nil
}

// TODO: doc -2024-10-31
type pageable func(cursor string) (next string, err error)

// TODO: doc -2024-10-31
func (s *Shredder) pager(fn pageable) (err error) {
	cursor := ""
	for {
		cursor, err = fn(cursor)
		if err != nil {
			return err
		}
		if cursor == "" {
			// Done - no more items to process.
			return nil
		}
		// Sleep for a bit to avoid rate limiting.
		time.Sleep(s.cfg.Sleep)
	}
}

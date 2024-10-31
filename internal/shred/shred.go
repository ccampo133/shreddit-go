package shred

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ccampo133/shreddit-go/internal/reddit"
)

const (
	// TODO: doc -ccampo 2024-10-30
	DefaultReplacementComment = "[deleted]"
)

// TODO: doc -ccampo 2024-10-30
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
	ReplacementComment string
}

// TODO: doc -ccampo 2024-10-30
type Shredder struct {
	client *reddit.Client
	cfg    Config
}

// TODO: doc -ccampo 2024-10-30
func NewShredder(client *reddit.Client, cfg Config) *Shredder {
	if cfg.Before.IsZero() {
		cfg.Before = time.Now()
	}
	if cfg.ReplacementComment == "" {
		cfg.ReplacementComment = DefaultReplacementComment
	}
	return &Shredder{client: client, cfg: cfg}
}

// TODO: doc -ccampo 2024-10-30
func (s *Shredder) Shred(ctx context.Context) error {
	// Comments
	if !s.cfg.SkipComments {
		if err := pager(ctx, s.shredComments); err != nil {
			return fmt.Errorf("error shredding comments: %w", err)
		}
	}
	// Posts
	if !s.cfg.SkipPosts {
		if err := pager(ctx, s.shredPosts); err != nil {
			return fmt.Errorf("error shredding posts: %w", err)
		}
	}
	// Saved comments
	if !s.cfg.SkipSavedComments {
		if err := pager(ctx, s.shredSavedComments); err != nil {
			return fmt.Errorf("error shredding saved comments: %w", err)
		}
	}
	// Saved posts
	if !s.cfg.SkipSavedPosts {
		if err := pager(ctx, s.shredSavedPosts); err != nil {
			return fmt.Errorf("error shredding saved posts: %w", err)
		}
	}
	return nil
}

// TODO: doc -ccampo 2024-10-30
func (s *Shredder) shredComments(ctx context.Context, after string) (string, error) {
	res, err := s.client.GetComments(ctx, s.cfg.Username, after)
	if err != nil {
		return "", fmt.Errorf("error getting comments: %w", err)
	}
	for _, comment := range res.Data.Children {
		// Skip comments younger than the cutoff time.
		if comment.Data.CreatedUTC.Time.After(s.cfg.Before) {
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
			slog.Info("Would shred comment", "permalink", comment.Data.Permalink)
			continue
		}
		// Edit the comment
		if err := s.client.EditComment(ctx, comment.Data.ID, s.cfg.ReplacementComment); err != nil {
			return "", fmt.Errorf("error editing comment: %w", err)
		}
		if !s.cfg.EditOnly {
			// Delete the comment
			if err := s.client.DeleteComment(ctx, comment.Data.ID); err != nil {
				return "", fmt.Errorf("error deleting comment: %w", err)
			}
		}
		slog.Info("Successfully shredded comment", "permalink", comment.Data.Permalink)
	}
	return res.Data.After, nil
}

func (s *Shredder) shredPosts(ctx context.Context, after string) (string, error) {
	// TODO: implement -ccampo 2024-10-30
	return "", nil
}

func (s *Shredder) shredSavedComments(ctx context.Context, after string) (string, error) {
	// TODO: implement -ccampo 2024-10-30
	return "", nil
}

func (s *Shredder) shredSavedPosts(ctx context.Context, after string) (string, error) {
	// TODO: implement -ccampo 2024-10-30
	return "", nil
}

type pageable func(ctx context.Context, cursor string) (string, error)

func pager(ctx context.Context, fn pageable) (err error) {
	after := ""
	for {
		after, err = fn(ctx, after)
		if err != nil {
			return err
		}
		if after == "" {
			// Done - no more items to process.
			return nil
		}
		// Sleep for a bit to avoid rate limiting.
		// TODO: make sleep configurable -ccampo 2024-10-31
		time.Sleep(2 * time.Second)
	}
}

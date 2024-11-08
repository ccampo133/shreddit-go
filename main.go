package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alecthomas/kong"
	"github.com/ccampo133/shreddit-go/internal/reddit"
	"github.com/ccampo133/shreddit-go/internal/shred"
)

type CLI struct {
	Username           string           `help:"Reddit username." short:"u" required:"" env:"SHREDDIT_USERNAME"`
	Password           string           `help:"Reddit password." short:"p" required:"" env:"SHREDDIT_PASSWORD"`
	ClientID           string           `help:"Reddit client ID." required:"" env:"SHREDDIT_CLIENT_ID"`
	ClientSecret       string           `help:"Reddit client secret." required:"" env:"SHREDDIT_CLIENT_SECRET"`
	DryRun             bool             `help:"Don't actually remove anything - just log what would be removed." env:"SHREDDIT_DRY_RUN"`
	ThingTypes         []string         `help:"Thing types to remove. Possible values: posts, comments, friends, saved-posts, saved-comments" `
	Before             time.Time        `help:"Remove things before this date." env:"SHREDDIT_BEFORE"`
	MaxDays            *int             `help:"Remove things older than this many days. Doesn't apply if using 'before'." env:"SHREDDIT_MAX_DAYS"`
	MaxScore           *int             `help:"Remove things with a karma score less than this." env:"SHREDDIT_MAX_SCORE"`
	ReplacementComment string           `help:"Comment to replace removed comments with." short:"r" default:"[deleted]" env:"SHREDDIT_REPLACEMENT_COMMENT"`
	UserAgent          string           `help:"Reddit user agent." default:"shreddit-go" env:"SHREDDIT_USER_AGENT"`
	GdprExportDir      string           `help:"The path of the directory of the unzipped GDPR export data. If set, will use the GDPR export data instead of Reddit's APIs for discovering your data." env:"SHREDDIT_GDPR_EXPORT_DIR"`
	EditOnly           bool             `help:"Only edit comments, don't remove them." env:"SHREDDIT_EDIT_ONLY"`
	Sleep              time.Duration    `help:"Time to sleep between requests." env:"SHREDDIT_SLEEP"`
	Version            kong.VersionFlag `name:"version" short:"v" help:"Print version information and quit"`
}

var (
	// version is the application version. It is intended to be set at compile
	// time via the linker (e.g. -ldflags="-X main.version=...").
	version = "dev"
)

func main() {
	cli := CLI{}
	ctx := kong.Parse(
		&cli,
		kong.Name("shreddit"),
		kong.Description("Overwrite and delete your Reddit account history."),
		kong.UsageOnError(),
		kong.ConfigureHelp(
			kong.HelpOptions{
				Compact: true,
			},
		),
		kong.Vars{
			"version": version,
		},
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

func (cli *CLI) Run() error {
	ctx := context.Background()
	redditCfg := reddit.Config{
		ClientID:     cli.ClientID,
		ClientSecret: cli.ClientSecret,
		Username:     cli.Username,
		Password:     cli.Password,
		UserAgent:    cli.UserAgent,
	}
	client, err := reddit.NewClient(ctx, redditCfg)
	if err != nil {
		return fmt.Errorf("error creating Reddit client: %w", err)
	}
	// TODO: check thing types to determine skip bools
	cfg := shred.Config{
		Username:           cli.Username,
		DryRun:             cli.DryRun,
		EditOnly:           cli.EditOnly,
		Before:             cli.Before,
		MaxScore:           cli.MaxScore,
		MaxDays:            cli.MaxDays,
		ReplacementComment: cli.ReplacementComment,
		Sleep:              cli.Sleep,
		// TODO: skip comments/posts/saved
	}
	shredder := shred.NewShredder(client, cfg)
	if err := shredder.Shred(); err != nil {
		return fmt.Errorf("error shredding: %w", err)
	}
	return nil
}

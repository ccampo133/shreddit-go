package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Username           string           `help:"Reddit username." short:"u" required:"" env:"SHREDDIT_USERNAME"`
	Password           string           `help:"Reddit password." short:"p" required:"" env:"SHREDDIT_PASSWORD"`
	ClientID           string           `help:"Reddit client ID." required:"" env:"SHREDDIT_CLIENT_ID"`
	ClientSecret       string           `help:"Reddit client secret." required:"" env:"SHREDDIT_CLIENT_SECRET"`
	DryRun             bool             `help:"Don't actually remove anything - just log what would be removed." env:"SHREDDIT_DRY_RUN"`
	ThingTypes         []string         `help:"Thing types to remove. Possible values: posts, comments, friends, saved-posts, saved-comments" `
	Before             time.Time        `help:"Remove things before this date." env:"SHREDDIT_BEFORE"`
	MaxScore           int              `help:"Remove things with a karma score less than this." env:"SHREDDIT_MAX_SCORE"`
	ReplacementComment string           `help:"Comment to replace removed comments with." short:"r" default:"[deleted]" env:"SHREDDIT_REPLACEMENT_COMMENT"`
	UserAgent          string           `help:"Reddit user agent." default:"shreddit-go" env:"SHREDDIT_USER_AGENT"`
	GdprExportDir      string           `help:"The path of the directory of the unzipped GDPR export data. If set, will use the GDPR export data instead of Reddit's APIs for discovering your data." env:"SHREDDIT_GDPR_EXPORT_DIR"`
	EditOnly           bool             `help:"Only edit comments, don't remove them." env:"SHREDDIT_EDIT_ONLY"`
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
	fmt.Println("Running shreddit")
	return nil
}

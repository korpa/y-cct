package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/korpa/y-cct/commands/replacer"
	"github.com/korpa/y-cct/commands/retention"
	"github.com/korpa/y-cct/global/commands/version"
	"github.com/korpa/y-cct/global/logging"
	"github.com/spf13/cobra"
)

var (
	Version       = "v0.0.0-dev"
	Goos          = "unknown"
	Goarch        = "unknown"
	Commit        = "unknown"
	Date          = "0000-00-00 00:00:00"
	BaseName      = "y"
	MarketingSlug = "y - Composite Command Tree"
	Description   = "This is a Composite Command Tree for Go Cobra"
)

var Cmd = &cobra.Command{
	Use:   BaseName,
	Short: BaseName + " commands",
}

func main() {
	versionInfo := version.VersionInfo{
		Version:       Version,
		Goos:          Goos,
		Goarch:        Goarch,
		Commit:        Commit,
		Date:          Date,
		BaseName:      BaseName,
		MarketingSlug: MarketingSlug,
		Description:   Description,
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "versionInfo", versionInfo)

	logger := logging.GetLogger()
	slog.SetDefault(logger)

	Cmd.AddCommand(replacer.Cmd)
	Cmd.AddCommand(retention.Cmd)
	Cmd.AddCommand(version.Cmd)

	if err := fang.Execute(
		ctx,
		Cmd,
		//fang.WithNotifySignal(os.Interrupt, os.Kill),
		fang.WithoutVersion(),
	); err != nil {
		os.Exit(1)
	}

}

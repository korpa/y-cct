package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/korpa/y-cct/commands/retention"
	"github.com/korpa/y-cct/global/commands/version"
	"github.com/korpa/y-cct/global/logging"
)

var (
	Version       = "v0.0.0-dev"
	Goos          = "unknown"
	Goarch        = "unknown"
	Commit        = "unknown"
	Date          = "0000-00-00 00:00:00"
	BaseName      = "y-retention"
	MarketingSlug = "y-retention - part of y-cct"
	Description   = ""
)

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

	if err := fang.Execute(
		context.Background(),
		retention.Cmd,
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	); err != nil {
		slog.Error(fmt.Sprint(err))
		os.Exit(1)
	}

}

func init() {
	retention.Cmd.AddCommand(version.Cmd)
}

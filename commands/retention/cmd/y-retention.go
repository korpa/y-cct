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
	Version  = "0.0.0_dev"
	Os       = "os_unknown"
	Arch     = "arch_unknown"
	Commit   = "commit_unknown"
	BaseName = "retention"
	Date     = "date_unknown"
)

func main() {

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

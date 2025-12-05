package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/korpa/y-cct/commands/replacer"
	"github.com/korpa/y-cct/global/commands/version"
	"github.com/korpa/y-cct/global/logging"
)

var (
	Version  = "0.0.0_dev"
	Os       = "os_unknown"
	Arch     = "arch_unknown"
	Commit   = "commit_unknown"
	BaseName = "replacer"
	Date     = "date_unknown"
)

func main() {
	logger := logging.GetLogger()
	slog.SetDefault(logger)

	if err := fang.Execute(
		context.Background(),
		replacer.Cmd,
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	); err != nil {
		slog.Error(fmt.Sprint(err))
		os.Exit(1)
	}
}

func init() {
	replacer.Cmd.AddCommand(version.Cmd)
}

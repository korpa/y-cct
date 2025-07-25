package main

import (
	"log/slog"

	"github.com/korpa/y-cct/commands/retention"
	"github.com/korpa/y-cct/global/commands/version"
	"github.com/korpa/y-cct/global/logging"
	"github.com/spf13/cobra"
)

var (
	Version  = "0.0.0_dev"
	Os       = "os_unknown"
	Arch     = "arch_unknown"
	Commit   = "commit_unknown"
	BaseName = "y"
	Date     = "date_unknown"
)

var Cmd = &cobra.Command{
	Use:   "y",
	Short: "y commands",
}

func main() {

	logger := logging.GetLogger()
	slog.SetDefault(logger)

	Cmd.AddCommand(retention.Cmd)
	Cmd.AddCommand(version.Cmd)

	Cmd.Execute()
}

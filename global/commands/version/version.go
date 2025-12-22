package version

import (
	"fmt"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type VersionInfo struct {
	Version       string
	Goos          string
	Goarch        string
	Commit        string
	Date          string
	BaseName      string
	MarketingSlug string
	Description   string
}

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number`,
	Run: func(cmd *cobra.Command, args []string) {

		ctx := cmd.Context()

		versionInfo := ctx.Value("versionInfo").(VersionInfo)

		if term.IsTerminal(int(os.Stdout.Fd())) {

			myFigure := figure.NewFigure(versionInfo.BaseName, "", true)
			myFigure.Print()

			fmt.Println()
			fmt.Println(versionInfo.MarketingSlug)
			fmt.Println()
			fmt.Println(versionInfo.Description)
			fmt.Println()
			fmt.Println(versionInfo.BaseName + " version " + versionInfo.Version + "-" + versionInfo.Commit + " " + versionInfo.Goos + "/" + versionInfo.Goarch + " - " + versionInfo.Date)
			fmt.Println()
		} else {
			fmt.Println(versionInfo.Version)
		}

	},
}

func init() {

}

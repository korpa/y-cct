package version

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version       = "developmentVERSION"
	Os            = "unknownOS"
	Arch          = "unknownARCH"
	Commit        = "unknownCOMMIT"
	BaseName      = "y"
	MarketingSlug = "y - Composite Command Tree"
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("short") {
			fmt.Println(Version + "-" + Commit + "-" + Os + "-" + Arch)
		} else {
			myFigure := figure.NewFigure(BaseName, "", true)
			myFigure.Print()
			fmt.Println()
			fmt.Println(MarketingSlug)
			fmt.Println()
			fmt.Println(`This is a Composite Command Tree for Go Cobra`)
			fmt.Println()
			fmt.Println(BaseName + " version " + Version + "-" + Commit + "-" + Os + "-" + Arch)
			fmt.Println()
		}
	},
}

func init() {
	Cmd.Flags().BoolP("short", "s", false, "show short version message")
	viper.BindPFlag("short", Cmd.Flags().Lookup("short"))
}

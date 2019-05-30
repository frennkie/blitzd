package main

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/blitzd"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Version: blitzd.BuildVersion,
	Use:     "blitzd",
	Short:   "RaspiBlitz Daemon",
	Long: `A service that retrieves and caches details about your RaspiBlitz.
More info at: https://github.com/frennkie/blitzd`,
	Run: func(cmd *cobra.Command, args []string) {
		blitzd.Init()
	},
}

var demoCmd = &cobra.Command{
	Version: blitzd.BuildVersion,
	Use:     "version",
	Short:   "Show detailed version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("BuildTime:", blitzd.BuildTime)
		fmt.Println("BuildVersion:", blitzd.BuildVersion)
	},
}

func main() {
	cobra.OnInitialize(config.InitConfig)

	RootCmd.PersistentFlags().StringVar(&config.BlitzdDir, "dir",
		config.DefaultBlitzdDir, "blitzd home directory (default is $HOME/.blitzd")

	RootCmd.AddCommand(demoCmd)

	_ = RootCmd.Execute()

}

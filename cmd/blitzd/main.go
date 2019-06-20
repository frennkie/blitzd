package main

import (
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/pkg/cmd/blitzd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.OnInitialize(config.InitConfig)

	blitzd.RootCmd.PersistentFlags().StringVar(&config.BlitzdDir, "dir",
		config.DefaultBlitzdDir, "blitzd home directory (default is $HOME/.blitzd")

	blitzd.RootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v",
		true, "print more log messages (default: true)")

	blitzd.RootCmd.AddCommand(blitzd.DemoCmd)
	blitzd.RootCmd.AddCommand(blitzd.GenCertCmd)
	blitzd.RootCmd.AddCommand(blitzd.GraceCmd)

	_ = blitzd.RootCmd.Execute()

}

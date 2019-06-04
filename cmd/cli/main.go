package main

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/pkg/cmd/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultRPCHostPort = fmt.Sprintf("localhost:%d", config.DefaultRPCPort)
)

func main() {
	cobra.OnInitialize(config.InitConfig)

	cli.RootCmd.PersistentFlags().StringVar(&config.BlitzdDir, "dir",
		config.DefaultBlitzdDir, "blitzd home directory (default is $HOME/.blitzd")

	cli.RootCmd.PersistentFlags().StringP("rpcHostPort", "H", defaultRPCHostPort, "Host and Port to connect to")
	_ = viper.BindPFlag("rpcHostPort", cli.RootCmd.PersistentFlags().Lookup("rpcHostPort"))
	viper.SetDefault("rpcHostPort", defaultRPCHostPort)

	cli.RootCmd.AddCommand(cli.CmdAll)
	cli.RootCmd.AddCommand(cli.CmdArch)
	cli.RootCmd.AddCommand(cli.CmdUptime)
	cli.RootCmd.AddCommand(cli.CmdTimes)
	cli.RootCmd.AddCommand(cli.CmdEcho)
	cli.RootCmd.AddCommand(cli.CmdApi)

	_ = cli.RootCmd.Execute()

}

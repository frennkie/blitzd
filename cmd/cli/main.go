package main

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/pkg/cmd/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cobra.OnInitialize(config.InitConfig)

	cli.RootCmd.PersistentFlags().StringVarP(&config.BlitzdDir, "dir", "D",
		config.DefaultBlitzdDir, "blitzd home directory (default is $HOME/.blitzd")

	cli.RootCmd.PersistentFlags().StringVarP(&config.RpcHostPort,
		"rpcHostPort", "H", fmt.Sprintf("localhost:%d", config.DefaultRPCPort),
		"Host and Port to connect to")
	_ = viper.BindPFlag("rpcHostPort", cli.RootCmd.PersistentFlags().Lookup("rpcHostPort"))

	cli.RootCmd.AddCommand(cli.CmdAll)
	cli.RootCmd.AddCommand(cli.CmdArch)
	cli.RootCmd.AddCommand(cli.CmdUptime)
	cli.RootCmd.AddCommand(cli.CmdTimes)
	cli.RootCmd.AddCommand(cli.CmdEcho)
	cli.RootCmd.AddCommand(cli.CmdApi)

	_ = cli.RootCmd.Execute()

}

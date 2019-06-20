package main

import (
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/pkg/cmd/cli"
	"github.com/spf13/cobra"
)

func main() {
	cobra.OnInitialize(config.InitConfig)

	cli.Init()

}

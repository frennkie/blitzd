package cmd

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	rootCmd.AddCommand(cmdUptime)
}

var cmdUptime = &cobra.Command{
	Use:   "uptime",
	Short: "uptime",
	Long:  `System uptime`,
	Run: func(cmd *cobra.Command, args []string) {
		hostPort := viper.GetString("restHostPort")
		url := "http://" + hostPort + data.APIv1

		myCache := new(data.Cache)

		if err := util.GetJson(url, myCache); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		println(myCache.Uptime.Value)
	},
}

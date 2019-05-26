package cmd

import (
	"github.com/frennkie/blitzinfod/internal/data"
	"github.com/frennkie/blitzinfod/internal/jsonclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(cmdAll)
}

var cmdAll = &cobra.Command{
	Use:   "all",
	Short: "Print full json blob",
	Long:  `Print full json blob`,
	Run: func(cmd *cobra.Command, args []string) {
		hostPort := viper.GetString("restHostPort")
		url := "http://" + hostPort + "/api/"

		myCache := new(data.Cache)
		_ = jsonclient.GetJson(url, myCache)
		println(myCache.Arch.Value)
		println(myCache.OperatingSystem.Value)
	},
}

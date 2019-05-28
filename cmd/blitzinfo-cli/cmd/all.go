package cmd

import (
	"fmt"
	"github.com/frennkie/blitzinfod/internal/data"
	"github.com/frennkie/blitzinfod/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
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
		url := "http://" + hostPort + data.APIv1

		myCache := new(data.Cache)

		if err := utils.GetJson(url, myCache); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		println("Arch: " + myCache.Arch.Value)
		println("OS: " + myCache.OperatingSystem.Value)
		println("Foo: " + myCache.Foo.Value)
		println("Uptime: " + myCache.Uptime.Value)
	},
}

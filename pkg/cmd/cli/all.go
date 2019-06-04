package cli

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var CmdAll = &cobra.Command{
	Use:   "all",
	Short: "Print full json blob",
	Long:  `Print full json blob`,
	Run: func(cmd *cobra.Command, args []string) {
		hostPort := viper.GetString("rpcHostPort")
		url := "http://" + hostPort + data.APIv1

		myCache := new(data.CacheOld)

		if err := util.GetJson(url, myCache); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		println("Arch: " + myCache.Arch.Value)
		println("OS: " + myCache.OperatingSystem.Value)
		println("Foo: " + myCache.Foo.Value)
		println("Uptime: " + myCache.Uptime.Value)
	},
}

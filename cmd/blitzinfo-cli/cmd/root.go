package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const (
	defaultRESTPort     = "7080"
	defaultRESTHostPort = "localhost:" + defaultRESTPort

	//defaultRPCPort          = "39735"
	//defaultRPCHostPort      = "localhost:" + defaultRPCPort
)

var rootCmd = &cobra.Command{
	Version: "1.2.3",
	Use:     "blitzinfo-cli",
	Short:   "blitzinfo-cli is the CLI for model",
	Long: `An easy way to access data from model.
                More info at: https://github.com/frennkie/model`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("Stuff")
		fmt.Println("default RESTHostPort: " + defaultRESTHostPort)
		fmt.Println("given RESTHostPort: " + viper.GetString("restHostPort"))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("restHostPort", "H", defaultRESTHostPort, "Host and Port to connect to")
	_ = viper.BindPFlag("restHostPort", rootCmd.PersistentFlags().Lookup("restHostPort"))
	viper.SetDefault("restHostPort", defaultRESTHostPort)
}

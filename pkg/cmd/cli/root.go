package cli

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/pkg/cmd/blitzd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultRPCPort = 39735
)

var jsonFlag bool
var formattedFlag bool

var rootCmd = &cobra.Command{
	Version: blitzd.BuildVersion + " (built: " + blitzd.BuildTime + ")",
	Use:     "blitz-cli",
	Short:   "blitz-cli is the CLI for blitzd",
	Long: `An easy way to access data from blitzd.
                More info at: https://github.com/frennkie/blitzd`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		log.WithFields(log.Fields{"rpcHostPort": config.RpcHostPort}).Debug("cli config")
	},
}

func Init() {

	rootCmd.PersistentFlags().StringVarP(&config.BlitzdDir, "dir", "D",
		config.DefaultBlitzdDir, "blitzd home directory")

	rootCmd.PersistentFlags().StringVarP(&config.RpcHostPort,
		"rpcHostPort", "H", fmt.Sprintf("localhost:%d", DefaultRPCPort),
		"Host and Port to connect to")
	_ = viper.BindPFlag("rpcHostPort", rootCmd.PersistentFlags().Lookup("rpcHostPort"))

	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v",
		false, "print debug log messages")

	rootCmd.PersistentFlags().BoolVarP(&config.Trace, "trace", "t",
		false, "print all (also debug and trace) log messages")
	_ = rootCmd.PersistentFlags().MarkHidden("trace")

	// ToDo(frennkie) remove these two
	rootCmd.AddCommand(cmdHello)
	rootCmd.AddCommand(cmdHelloWorld)

	rootCmd.AddCommand(cmdGet)
	cmdGet.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output as JSON")
	cmdGet.Flags().BoolVarP(&formattedFlag, "formatted", "f", false, "Output as formatted value")
	cmdGet.AddCommand(cmdGetAll)

	rootCmd.AddCommand(cmdFoo5)

	rootCmd.AddCommand(cmdShutdown)

	_ = rootCmd.Execute()

}

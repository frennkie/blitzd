package main

import (
	"bufio"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/frennkie/blitzd/internal/api"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
	_ "github.com/frennkie/blitzd/internal/metric/lnd"
	_ "github.com/frennkie/blitzd/internal/metric/network"
	_ "github.com/frennkie/blitzd/internal/metric/system"
	"github.com/frennkie/blitzd/internal/serve"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	defaultCfgFile      = "/etc/blitzd.toml"
	defaultCfgFileWin32 = "C:\\blitzd.toml"

	defaultRESTPort     = "7080"
	defaultRESTHostPort = "localhost:" + defaultRESTPort

	//defaultRPCPort     = "39735"
	//defaultRPCHostPort = "localhost:" + defaultRPCPort

	Normal = "normal"
	Red    = "red"
	Green  = "green"
	Yellow = "yellow"
	Purple = "purple"
)

var (
	cfgFile string

	buildVersion = "0.8.15"
	buildTime    = "0"
)

func UpdateFileBar() {
	title := "file-bar"
	absFilePath := "/tmp/foo"

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		log.Printf("file does not exist - skipping: %s", absFilePath)
		return
	}

	log.Printf("starting goroutine: %s (%s)", title, absFilePath)

	UpdateFileBarFunc(title, absFilePath)
	go util.FileWatcher(title, absFilePath, UpdateFileBarFunc)
}

func UpdateFileBarFunc(title string, absFilePath string) {
	log.Printf("event-based update: %s (%s)", title, absFilePath)
	m := data.NewMetricEventBased(title)

	m.Value = fmt.Sprintf("%s", "foobar")

	metric.MetricsMux.Lock()
	metric.Metrics.FileBar = m
	metric.MetricsMux.Unlock()

}

// TODO replace with /etc/issue
func UpdateLsbRelease() {
	title := "lsb-release"
	absFilePath := "/etc/lsb-release"

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		log.Printf("file does not exist - skipping: %s", absFilePath)
		return
	}

	log.Printf("starting goroutine: %s (%s)", title, absFilePath)

	UpdateLsbReleaseFunc(title, absFilePath)
	go util.FileWatcher(title, absFilePath, UpdateLsbReleaseFunc)
}

func UpdateLsbReleaseFunc(title string, absFilePath string) {
	log.Printf("event-based update: %s (%s)", title, absFilePath)
	m := data.NewMetricEventBased(title)

	file, err := os.Open(absFilePath)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	_ = file.Close()

	tmp := txtlines[len(txtlines)-1]
	tmp2 := strings.Split(tmp, "=")[1]
	tmp3 := strings.Replace(tmp2, "\"", "", -1)
	m.Value = tmp3

	metric.MetricsMux.Lock()
	metric.Metrics.LsbRelease = m
	metric.MetricsMux.Unlock()

}

func main() {
	var rootCmd = &cobra.Command{
		Version: buildVersion,
		Use:     "blitzd",
		Short:   "RaspiBlitz Daemon",
		Long: `A service that retrieves and caches details about your RaspiBlitz.
More info at: https://github.com/frennkie/blitzd`,
		Run: func(cmd *cobra.Command, args []string) {
			blitzd()
		},
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.blitzd.cfg")
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().StringP("restHostPort", "R", defaultRESTHostPort, "REST: Listen on Host:Port")
	_ = viper.BindPFlag("restHostPort", rootCmd.PersistentFlags().Lookup("restHostPort"))
	viper.SetDefault("restHostPort", defaultRESTHostPort)

	//rootCmd.PersistentFlags().StringP("rpcHostPort", "H", defaultRPCHostPort, "RPC: Listen on Host:Port")
	//_ = viper.BindPFlag("rpcHostPort", rootCmd.PersistentFlags().Lookup("rpcHostPort"))
	//viper.SetDefault("rpcHostPort", defaultRPCHostPort)

	_ = rootCmd.Execute()
}

func initConfig() {
	// If config is specified by flag then ONLY read that file.
	// Otherwise read default (/etc/blitzd.toml) and - if it exists - merge any
	// settings from file "blitzd.toml" in home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}

	} else {
		// First read default config from /etc (or C:\ on Win32)
		if runtime.GOOS == "windows" {
			viper.SetConfigFile(filepath.FromSlash(defaultCfgFileWin32))
		} else {
			viper.SetConfigFile(filepath.FromSlash(defaultCfgFile))
		}

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		customConfig := filepath.FromSlash(home + "/blitzd.toml")
		if _, err := os.Stat(customConfig); os.IsNotExist(err) {
			log.Printf("custom config file does not exist - skipping: %s", customConfig)
			return
		}

		viper.SetConfigFile(filepath.FromSlash(customConfig))
		if err := viper.MergeInConfig(); err != nil {
			fmt.Println("Can't read config for merge:", err)
			os.Exit(1)
		}

		log.Printf("Merged config file: %s", customConfig)

	}

}

func blitzd() {
	log.Printf("Starting version: %s, built at %s", buildVersion, buildTime)

	if runtime.GOOS != "windows" {
		go UpdateLsbRelease()
		go UpdateFileBar()
	}

	box := rice.MustFindBox("../../web/")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(box.HTTPBox())))

	http.HandleFunc("/", serve.Root)
	// ToDo fix metrics
	http.HandleFunc("/info/", serve.Info(&metric.Metrics))

	// ToDo fix.. every sub url matches
	http.HandleFunc(data.APIv1, api.All())
	http.HandleFunc(data.APIv1+"config/", api.Config())
	http.HandleFunc(data.APIv1+"lnd/", api.Lnd())
	http.HandleFunc(data.APIv1+"system/", api.System())

	RESTHostPort := viper.GetString("RESTHostPort")
	log.Printf("REST: Listening on host: http://%s", RESTHostPort)

	//rpcHostPort := viper.GetString("rpcHostPort")
	//log.Printf("RPC: Listening on host: gRPC://%s", rpcHostPort)

	// now ListenAndServer
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", RESTHostPort), nil))

}

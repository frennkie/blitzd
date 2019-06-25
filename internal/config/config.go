package config

import (
	"bytes"
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	defaultBlitzdDirName    = "blitzd"
	defaultBlitzdConfigName = "config.toml"

	defaultCfgFile      = "/etc/blitzd.toml"
	defaultCfgFileWin32 = "C:\\blitzd.toml" // Win32 mostly used for development

	// default admin user is: admin
	// default password for admin is: changeme
	// use htpasswd (part of apache2-utils) to create a new password:
	//   htpasswd -n -B admin
	defaultAdminUsername = "admin"
	defaultAdminPassword = "$2y$05$nNUGiiHDDric6W/Zml05Ku0Ij04mf62NTd/JRWQya8uxLpoGR3yJS"

	defaultAlias = "MyBlitz42"

	defaultEnvPrefix = "BLITZD"
)

var (
	C *Config

	DefaultBlitzdDir = btcutil.AppDataDir(defaultBlitzdDirName, false)
	BlitzdDir        string
	RpcHostPort      string
	Verbose          bool
	Trace            bool

	// ToDo check
	// maxMsgRecvSize is the largest message our client will receive. We
	// set this to ~50Mb atm.
	//maxMsgRecvSize = grpc.MaxCallRecvMsgSize(1 * 1024 * 1024 * 50)

)

//
type Config struct {
	// ToDo(frennkie) check BlitzDir
	BlitzdDir string `mapstructure:"blitzd_dir" toml:"blitzd_dir"`

	// ToDo(frennkie) check these 2 (are set in config.InitConfig)
	CustomConfigPath  string `mapstructure:"custom_config_path" toml:"custom_config_path"`
	DefaultConfigPath string `mapstructure:"default_config_path" toml:"default_config_path"`

	Alias   string  `mapstructure:"alias" toml:"alias"`
	Admin   Admin   `mapstructure:"admin" toml:"admin"`
	Module  Module  `mapstructure:"module" toml:"module"`
	Client  Client  `mapstructure:"client" toml:"client"`
	Server  Server  `mapstructure:"server" toml:"server"`
	Service Service `mapstructure:"service" toml:"service"`
}

type Admin struct {
	Password string `mapstructure:"password" toml:"password"`
	Username string `mapstructure:"username" toml:"username"`
}

type Module struct {
	Bitcoind   Bitcoind   `mapstructure:"bitcoind" toml:"bitcoind"`
	Lnd        Lnd        `mapstructure:"lnd" toml:"lnd"`
	Network    Network    `mapstructure:"network" toml:"network"`
	RaspiBlitz RaspiBlitz `mapstructure:"raspiblitz" toml:"raspiblitz"`
	System     System     `mapstructure:"system" toml:"system"`
}

type Bitcoind struct {
	Enabled     bool   `mapstructure:"enabled" toml:"enabled"`
	RpcAddress  string `mapstructure:"rpc_address" toml:"rpc_address"`
	RpcPassword string `mapstructure:"rpc_password" toml:"rpc_password"`
	RpcUser     string `mapstructure:"rpc_user" toml:"rpc_user"`
}

type Lnd struct {
	Enabled    bool   `mapstructure:"enabled" toml:"enabled"`
	Macaroon   string `mapstructure:"macaroon" toml:"macaroon"`
	RpcAddress string `mapstructure:"rpc_address" toml:"rpc_address"`
	TlsCert    string `mapstructure:"tls_cert" toml:"tls_cert"`
}

type Network struct {
	Enabled bool   `mapstructure:"enabled" toml:"enabled"`
	Nic     string `mapstructure:"nic" toml:"nic"`
}

type RaspiBlitz struct {
	Enabled bool   `mapstructure:"enabled" toml:"enabled"`
	Config  string `mapstructure:"config" toml:"config"`
}

type System struct {
	Enabled bool   `mapstructure:"enabled" toml:"enabled"`
	Mount1  string `mapstructure:"mount1" toml:"mount1"`
	Mount2  string `mapstructure:"mount2" toml:"mount2"`
}

type Client struct {
	CaCert  string `mapstructure:"ca_cert" toml:"ca_cert"`
	TlsCert string `mapstructure:"tls_cert" toml:"tls_cert"`
	TlsKey  string `mapstructure:"tls_key" toml:"tls_key"`
}

type Server struct {
	CaCert   string          `mapstructure:"ca_cert" toml:"ca_cert"`
	TlsCert  string          `mapstructure:"tls_cert" toml:"tls_cert"`
	TlsKey   string          `mapstructure:"tls_key" toml:"tls_key"`
	Http     ServerConfig    `mapstructure:"http" toml:"http"`
	Https    ServerConfig    `mapstructure:"https" toml:"https"`
	HttpsTor ServerConfigTor `mapstructure:"https_tor" toml:"https_tor"`
	Rest     ServerConfig    `mapstructure:"rest" toml:"rest"`
	RestTor  ServerConfigTor `mapstructure:"rest" toml:"rest_tor"`
	Rpc      ServerConfig    `mapstructure:"rpc" toml:"rpc"`
	RpcTor   ServerConfigTor `mapstructure:"rpc_tor" toml:"rpc_tor"`
}

type ServerConfig struct {
	Enabled       bool `mapstructure:"enabled" toml:"enabled"`
	LocalhostOnly bool `mapstructure:"localhost_only" toml:"localhost_only"`
	Port          int  `mapstructure:"port" toml:"port"`
}

type ServerConfigTor struct {
	Enabled bool `mapstructure:"enabled" toml:"enabled"`
	Port    int  `mapstructure:"port" toml:"port"`
}

type Service struct {
	Shutdown Shutdown `mapstructure:"shutdown" toml:"shutdown"`
}

type Shutdown struct {
	Enabled bool   `mapstructure:"enabled" toml:"enabled"`
	Script  string `mapstructure:"script" toml:"script"`
}

// set Default values
func NewConfig() *Config {
	log.Debug("Settings Defaults")
	return &Config{
		BlitzdDir:         BlitzdDir,
		CustomConfigPath:  "",
		DefaultConfigPath: "",
		Alias:             defaultAlias,
		Admin: Admin{
			Password: defaultAdminPassword,
			Username: defaultAdminUsername,
		},
		Module: Module{
			Bitcoind: Bitcoind{
				Enabled:     false,
				RpcAddress:  "localhost:8332",
				RpcPassword: "bitcoin_rpc_password",
				RpcUser:     "raspibolt",
			},
			Lnd: Lnd{
				Enabled:    false,
				Macaroon:   "/home/bitcoin/.lnd/data/chain/bitcoin/mainnet/readonly.macaroon",
				RpcAddress: "localhost:10009",
				TlsCert:    "/home/bitcoin/.lnd/tls.cert",
			},
			Network: Network{
				Enabled: false,
				Nic:     "eth0",
			},
			RaspiBlitz: RaspiBlitz{
				Enabled: false,
				Config:  "/mnt/hdd/raspiblitz.conf",
			},
			System: System{
				Enabled: true,
				Mount1:  "/",
				Mount2:  "/mnt/hdd/",
			},
		},
		Client: Client{
			CaCert:  filepath.Join(BlitzdDir, "blitzd_ca.crt"),
			TlsCert: filepath.Join(BlitzdDir, "blitzd_client.crt"),
			TlsKey:  filepath.Join(BlitzdDir, "blitzd_client.key"),
		},
		Server: Server{
			CaCert:  filepath.Join(BlitzdDir, "blitzd_ca.crt"),
			TlsCert: filepath.Join(BlitzdDir, "blitzd_server.crt"),
			TlsKey:  filepath.Join(BlitzdDir, "blitzd_server.key"),
			Http: ServerConfig{
				Enabled:       true,
				LocalhostOnly: true,
				Port:          39080,
			},
			Https: ServerConfig{
				Enabled:       true,
				LocalhostOnly: true,
				Port:          39443,
			},
			HttpsTor: ServerConfigTor{
				Enabled: false,
				Port:    39444,
			},
			Rest: ServerConfig{
				Enabled:       false,
				LocalhostOnly: true,
				Port:          39445,
			},
			RestTor: ServerConfigTor{
				Enabled: false,
				Port:    39446,
			},
			Rpc: ServerConfig{
				Enabled:       true,
				LocalhostOnly: true,
				Port:          39735,
			},
			RpcTor: ServerConfigTor{
				Enabled: false,
				Port:    39736,
			},
		},
		Service: Service{
			Shutdown: Shutdown{
				Enabled: false,
				Script:  "/home/admin/XXreboot.sh",
			},
		},
	}
}

func SetupLogger() {
	// setup logrus
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	// Default is to show "Info" and above.
	// Verbose enables "Debug".
	// Trace enables "Debug" and "Trace".
	if Trace {
		log.SetLevel(log.TraceLevel)
	} else {
		if Verbose {
			log.SetLevel(log.DebugLevel)
		}
	}
}

func InitConfig() {
	//viper := viper.New()

	SetupLogger()

	// First set default values.
	// Then read default (/etc/blitzd.toml|C:\blitzd.toml)
	// Then - if it exists - merge any settings from file "blitzd.toml" in home directory
	// Then load env
	// Finally activate the Config Watcher

	// set default values in viper.
	// Viper needs to know if a key exists in order to override it.
	// https://github.com/spf13/viper/issues/188
	b, _ := toml.Marshal(NewConfig())
	defaultConfig := bytes.NewReader(b)
	viper.SetConfigType("toml")

	_ = viper.MergeConfig(defaultConfig)
	// refresh configuration with all merged values
	_ = viper.Unmarshal(&C)

	// read default
	if runtime.GOOS == "windows" {
		viper.Set("default_config_path", defaultCfgFileWin32)
		viper.SetConfigFile(filepath.FromSlash(defaultCfgFileWin32))
	} else {
		viper.Set("default_config_path", defaultCfgFile)
		viper.SetConfigFile(filepath.FromSlash(defaultCfgFile))
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Can't read config:", err)
	}

	// check custom config
	if _, err := os.Stat(filepath.FromSlash(BlitzdDir)); os.IsNotExist(err) {
		// BlitzdDir does not exist
		err = os.MkdirAll(filepath.FromSlash(BlitzdDir), 0700)
		if err != nil {
			log.Warn("Dir does not exist and can't be created: ", BlitzdDir)
			log.Fatal("err: ", err)
		}
	} else {

	}

	customCfgPath := filepath.FromSlash(filepath.Join(BlitzdDir, defaultBlitzdConfigName))
	if _, err := os.Stat(customCfgPath); os.IsNotExist(err) {
		if Verbose {
			log.Debug("custom config file does not exist - skipping: %s", customCfgPath)
		}
	} else {
		viper.Set("custom_config_path", customCfgPath)
		viper.SetConfigFile(customCfgPath)
		if err := viper.MergeInConfig(); err != nil {
			log.Fatal("Can't read config for merge:", err)
		}

		log.Debug("Merged config file: %s", customCfgPath)
	}

	// load env
	viper.SetEnvPrefix(defaultEnvPrefix)
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// refresh configuration with all merged values
	_ = viper.Unmarshal(&C)

	// full parsed/merged config as bytes
	bs, err := toml.Marshal(C)
	if err != nil {
		log.Fatalf("unable to marshal config to TOML: %v", err)
	}
	if Verbose {
		fmt.Println(string(bs))
	}

	// ToDo(frennkie) remove this
	// store copy of parsed/merged config
	absPath := filepath.Join(BlitzdDir, "saved.toml")
	if err := ioutil.WriteFile(absPath, bs, 0600); err != nil {
		log.WithFields(log.Fields{"absPath": absPath}).
			Warn("unable to write copy of parsed/merged config")
	}

	// config watcher
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.WithFields(log.Fields{"file": e.Name}).Debug("Config file changed")
	})

}

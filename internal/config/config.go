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
	Tor     Tor     `mapstructure:"tor" toml:"tor"`
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
	Tls Tls `mapstructure:"tls" toml:"tls"`
}

type Rest struct {
	Enabled bool `mapstructure:"enabled" toml:"enabled"`
	Docs    bool `mapstructure:"docs" toml:"docs"`
}

type Server struct {
	Tls   Tls               `mapstructure:"tls" toml:"tls"`
	Http  ServerConfig      `mapstructure:"http" toml:"http"`
	Https ServerConfigHttps `mapstructure:"https" toml:"https"`
	Grpc  ServerConfig      `mapstructure:"grpc" toml:"grpc"`
}

type ServerConfig struct {
	Enabled       bool `mapstructure:"enabled" toml:"enabled"`
	LocalhostOnly bool `mapstructure:"localhost_only" toml:"localhost_only"`
	Port          int  `mapstructure:"port" toml:"port"`
}

type ServerConfigHttps struct {
	Enabled       bool `mapstructure:"enabled" toml:"enabled"`
	LocalhostOnly bool `mapstructure:"localhost_only" toml:"localhost_only"`
	Port          int  `mapstructure:"port" toml:"port"`
	Rest          Rest `mapstructure:"rest" toml:"rest"`
}

type Service struct {
	Shutdown Shutdown `mapstructure:"shutdown" toml:"shutdown"`
}

type Shutdown struct {
	Enabled bool   `mapstructure:"enabled" toml:"enabled"`
	Script  string `mapstructure:"script" toml:"script"`
}

type Tls struct {
	Ca   string `mapstructure:"ca" toml:"ca"`
	Cert string `mapstructure:"cert" toml:"cert"`
	Key  string `mapstructure:"key" toml:"key"`
}

type Tor struct {
	Enabled        bool       `mapstructure:"enabled" toml:"enabled"`
	ExePath        string     `mapstructure:"exe_path" toml:"exe_path"`
	DataDir        string     `mapstructure:"data_dir" toml:"data_dir"`
	Hostname       string     `mapstructure:"hostname" toml:"hostname"`
	SecretKeyPath  string     `mapstructure:"secret_key_path" toml:"secret_key_path"`
	PublicKeyPath  string     `mapstructure:"public_key_path" toml:"public_key_path"`
	ServiceVersion int        `mapstructure:"service_version" toml:"service_version"`
	Tls            Tls        `mapstructure:"tls" toml:"tls"`
	Https          TorService `mapstructure:"https" toml:"https"`
	Grpc           TorService `mapstructure:"grpc" toml:"grpc"`
}

type TorService struct {
	Enabled   bool `mapstructure:"enabled" toml:"enabled"`
	LocalPort int  `mapstructure:"local_port" toml:"local_port"`
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
			Tls: Tls{
				Ca:   filepath.Join(BlitzdDir, "blitzd_ca.crt"),
				Cert: filepath.Join(BlitzdDir, "blitzd_client.crt"),
				Key:  filepath.Join(BlitzdDir, "blitzd_client.key"),
			},
		},
		Server: Server{
			Tls: Tls{
				Ca:   filepath.Join(BlitzdDir, "blitzd_ca.crt"),
				Cert: filepath.Join(BlitzdDir, "blitzd_server.crt"),
				Key:  filepath.Join(BlitzdDir, "blitzd_server.key"),
			},
			Http: ServerConfig{
				Enabled:       true,
				LocalhostOnly: true,
				Port:          39080,
			},
			Https: ServerConfigHttps{
				Enabled:       true,
				LocalhostOnly: true,
				Port:          39443,
				Rest: Rest{
					Enabled: false,
					Docs:    false,
				},
			},
			Grpc: ServerConfig{
				Enabled:       true,
				LocalhostOnly: true,
				Port:          39735,
			},
		},
		Service: Service{
			Shutdown: Shutdown{
				Enabled: false,
				Script:  "/home/admin/XXreboot.sh",
			},
		},
		Tor: Tor{
			Enabled:        false,
			ExePath:        "/usr/sbin/tor",
			DataDir:        filepath.Join(BlitzdDir, "tor"),
			Hostname:       "",
			SecretKeyPath:  "",
			PublicKeyPath:  "",
			ServiceVersion: 3,
			Tls: Tls{
				Ca:   filepath.Join(BlitzdDir, "tor", "blitzd_tor_ca.crt"),
				Cert: filepath.Join(BlitzdDir, "tor", "blitzd_tor_server.crt"),
				Key:  filepath.Join(BlitzdDir, "tor", "blitzd_tor_server.key"),
			},
			Https: TorService{
				Enabled:   false,
				LocalPort: 39444,
			},
			Grpc: TorService{
				Enabled:   false,
				LocalPort: 39736,
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

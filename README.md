# Blitz Info Daemon


#### install go (1.12+)

#### add go bin dir to PATH (either one time or add to .zshrc)

```bash
export GOPATH="$HOME/work/go" # add GOPATH (golang)
export GOBIN="$GOPATH/bin"    # add GOBIN (binary)
export PATH="$PATH:$GOPATH"   # add GOPATH to PATH
```


#### make build

```bash
foo
```


### Resources

https://github.com/golang-standards/project-layout

https://stackoverflow.com/questions/52899535/using-logger-configs-in-multi-packages-best-practice-for-golang-productive




### Configuration

Configuration values can be set in the following ways in the given order (the last will *win*):

* `/etc/blitzd.toml`
* `$HOME/.blitzd/config.toml` (this path can be customized on CLI)
* `env` (system environment)
* command line flags (where available)

**Please note:** blitzd must be restarted after every configuration change!


Interface

Für Cache -> Set (der macht auch den Lock)

-> Atomic .. Singleton

### Additional

for `htpasswd`:

``` 
sudo apt-get install apache2-utils
```

## Dev/Building

### Dev/Build Requirements

sudo apt-get update
sudo apt-get install go-dep devscripts libdistro-info-perl dh-systemd protobuf-compiler


### Git Tag

git tag vX.Y[-rcZ] -m "Message.."


### Building Blocks

* Logging and Errors 
  * errors - https://github.com/juju/errors 
  * logrus - https://github.com/sirupsen/logrus
* CLI and Config Management
  * cobra (CLI) - https://github.com/spf13/cobra
  * viper (Config) - https://github.com/spf13/viper
* Cache/In-Memory Data Storage
  * go-cache - https://github.com/patrickmn/go-cache
* API -> Protocol Buffers (protobuf)
  * gRPC - https://google.golang.org/grpc
* Web
  * net/http
  * html/template **OR** plush  
  * vfsgen (embedding static asset) https://github.com/shurcooL/vfsgen
  * httpauth (HTTP Basic Authentication) https://github.com/goji/httpauth


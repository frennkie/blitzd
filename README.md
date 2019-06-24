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


## "Proposed Connectivity"

|    |                 | Welcome   | Info            | REST        | gRPC        |                                                                                  |
|----|-----------------|-----------|-----------------|-------------|-------------|----------------------------------------------------------------------------------|
| 1  | description     |           |                 |             |             | T is True, F is False, in Brackets () is default value                           |
| 2  | enabled         | T/F (T)   | T/F (F)         | T/F (F)     | T/F (T)     |                                                                                  |
| 3  | localhost_only  | T/F (T)   | T/F (T)         | T/F (T)     | T/F (T)     |                                                                                  |
| 4  | TCP Port        | 80*/39080 | 443*/39443      | 39735       | 39737       | 80 and 443 required privileges/additional setting (setcap)                       |
| 5  | HTTPS/TLS       | False     | True            | True        | True        |                                                                                  |
| 6  | Certificate     | N/A       | Yes             | Yes         | Yes         | either self-signed, signed by self-signed CA or trusted CA (e.g. LE)             |
| 7  | Authentication  | None      | Basic Auth      | Basic Auth  | Mutual TLS  | Basic Auth PW: in config file?! Mutual: Clients need trusted Cert                |
| 8  | Tor - enabled   | N/A       | T/F (F)         | T/F (F)     | T/F (F)     |                                                                                  |
| 9  | Tor - int Port  | N/A       | 39444           | 39736       | 39738       |                                                                                  |
| 10 | Tor - ext Port  | N/A       | 443/39443/39444 | 39735/39736 | 39737/39738 |                                                                                  |
| 11 | Tor - HTTPS/TLS | N/A       | Yes             | Yes         | Yes         |                                                                                  |
| 12 | Tor - Cert      | N/A       | Yes             | Yes         | Yes         | should definitely not be the same as above in line 6. So most likely self-signed |
| 13 | Tor - Auth      | N/A       | Basic Auth      | Basic Auth  | Mutual TLS  | same as line 7 ?!                                                                |
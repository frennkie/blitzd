# Blitz Info Daemon


#### install go (1.12+)

#### add go bin dir to PATH (either one time or add to .zshrc)

```bash
export GOPATH="$HOME/work/go" # add GOPATH (golang)
export GOBIN="$GOPATH/bin"    # add GOBIN (binary)
export PATH="$PATH:$GOPATH"   # add GOPATH to PATH
```


#### install rice (for embedding templates in binary)

```bash
go get github.com/GeertJohan/go.rice
go get github.com/GeertJohan/go.rice/rice
```


#### make build

```bash
foo
```


### Resources

https://github.com/golang-standards/project-layout


### Configuration

Configuration values can be set in the following ways in the given order (the last will *win*):

* `/etc/blitzd.toml`
* `$HOME/.blitzd/config.toml` (this path can be customized on CLI)
* `env` (system environment)
* command line flags (where available)

**Please note:** blitzd must be restarted after every configuration change!
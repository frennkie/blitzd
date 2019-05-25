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
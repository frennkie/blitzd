GOBUILD := go build -v
RICE := rice embed-go

build:
	$(RICE)
	$(GOBUILD) -o blitzinfod-amd64 server.go rice-box.go
	GOARCH=arm GOARM=7 $(GOBUILD) -o blitzinfod-armv7 server.go rice-box.go

xkcddir=./xkcdserver/cmd/xkcd
webdir=./webserver/cmd
XKCDTARGET=xkcd-server
WEBTARGET=web-server

all: webserver xkcdserver

webserver: deps
	@go build -o $(WEBTARGET) $(webdir)

xkcdserver: deps
	@go build -o $(XKCDTARGET) $(xkcddir)

deps:
	@go mod tidy

test:
	@go test -race -coverprofile=coverage.out ./xkcdserver/...
	@go tool cover -html=coverage.out -o coverage.html

lint:
	@golangci-lint run -v

sec:
	@trivy fs .
	@govulncheck ./...

e2e:
	@bash ./xkcdserver/e2e.sh

clean:
	@rm $(XKCDTARGET) $(WEBTARGET)

.PHONY: build clean deps

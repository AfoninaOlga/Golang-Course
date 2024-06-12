srcdir=./xkcdserver/cmd/xkcd
TARGET=xkcd-server

build: deps
	@go build -o $(TARGET) $(srcdir)

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
	@chmod +x e2e.sh
	@./e2e.sh

clean:
	@rm $(TARGET)

.PHONY: build clean deps

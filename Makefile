srcdir=./cmd/xkcd
TARGET=xkcd-server

build: deps
	go build -o $(TARGET) $(srcdir)

deps:
	go mod tidy

test:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run -v

sec:
	trivy fs .
	govulncheck ./...

clean:
	rm $(TARGET)

.PHONY: build clean deps

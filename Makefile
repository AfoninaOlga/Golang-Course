srcdir=./cmd/xkcd
TARGET=xkcd-server

build: deps
	go build -o $(TARGET) $(srcdir)

deps:
	go mod tidy

clean:
	rm $(TARGET)

.PHONY: build clean deps

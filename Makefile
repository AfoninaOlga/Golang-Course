srcdir=./cmd/xkcd
TARGET=xkcd

build: install
	go build -o $(TARGET) $(srcdir)

install:
	go mod tidy

clean:
	rm $(TARGET)

.PHONY: build clean install

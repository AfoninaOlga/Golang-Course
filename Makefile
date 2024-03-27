PROJECT_DIR=./stem
BIN_NAME=myapp

build:
	go build -o $(BIN_NAME) $(PROJECT_DIR)

clean:
	rm $(BIN_NAME)

default:
	build

.PHONY: build clean
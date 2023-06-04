BIN=smalte

build:
	@go build -ldflags="-extldflags=-static" -o $(BIN) main.go \
	&& strip $(BIN)
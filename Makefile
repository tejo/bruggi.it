.PHONY: all build build-arm serve clean

APP_NAME = bruggi
BIN_DIR = bin

all: build

build:
	@echo "Building for host architecture..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) main.go

build-arm:
	@echo "Building for Raspberry Pi (ARM64)..."
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=arm64 go build -o $(BIN_DIR)/$(APP_NAME)-arm64 main.go

serve:
	@echo "Running in development mode..."
	go run main.go -serve

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR) dist

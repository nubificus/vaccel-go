BIN_DIR = bin

PKG_CONFIG_PC_PATH := $(shell pkg-config --variable pc_path pkg-config)
PKG_CONFIG_ENV_PATH := $(value PKG_CONFIG_PATH)
export PKG_CONFIG_PATH := $(PKG_CONFIG_PC_PATH)$(if $(PKG_CONFIG_ENV_PATH),:$(PKG_CONFIG_ENV_PATH))

.PHONY: all prepare clean
all: noop classify detect exec nonser tf tflite torch

prepare:
	@go mod tidy
	@mkdir -p $(BIN_DIR)/

%: prepare examples/%/main.go
	go build -o $(BIN_DIR)/$@ examples/$@/main.go

clean:
	rm -rf $(BIN_DIR)

ROOT_DIR := $(abspath  $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))/)
BIN_PATH = ${ROOT_DIR}/bin/
OUTPUT = ${BIN_PATH}

build_bin: clean
	@echo "=====>\tTarget: build_bin"
	@go build -o ${OUTPUT}/funds-system ${ROOT_DIR}/main.go
	@cp ${ROOT_DIR}/config.yaml ${OUTPUT}
	@echo "output path: ${OUTPUT}"

clean:
	@echo "=====> \tTarget: clean"
	@echo "cleaning ${BIN_PATH}"
	@rm -rf ${BIN_PATH}
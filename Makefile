ROOT_DIR := $(abspath  $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))/)
CONFIG := ${ROOT_DIR}/config
OUTPUT := ${ROOT_DIR}/bin

build_service: 
	@echo "Target: build service"
	@echo "cleaning ${OUTPUT}"
	@rm -rf ${OUTPUT}
	@go build -o ${OUTPUT}/MarketBrain ${ROOT_DIR}/main.go
	@cp -R ${CONFIG} ${OUTPUT}
	@echo "output path: ${OUTPUT}"
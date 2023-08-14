ROOT_DIR := $(abspath  $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))/)
CMD := ${ROOT_DIR}/cmd
OUTPUT := ${ROOT_DIR}/bin

no_target:
	@echo "no target unspecialed"

build_funds_system: CMD := ${CMD}/funds_system
build_funds_system: OUTPUT := ${OUTPUT}/funds_system
build_funds_system: build_service

build_service: 
	@echo "Target: build service"
	@echo "cleaning ${OUTPUT}"
	@rm -rf ${OUTPUT}
	@go build -o ${OUTPUT}/main ${CMD}/main.go
	@cp ${CMD}/config.yaml ${OUTPUT}
	@echo "output path: ${OUTPUT}"
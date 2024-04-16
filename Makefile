BIN_DIR=_output/bin

# If tag not explicitly set in users default to the git sha.
TAG ?= v1.26.9-scheduler

.EXPORT_ALL_VARIABLES:

all: local

init:
	mkdir -p ${BIN_DIR}

local: init
	go build -o=${BIN_DIR}/tanjunchen-scheduler ./cmd/scheduler

build-linux: init
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o=${BIN_DIR}/tanjunchen-scheduler ./cmd/scheduler

image: build-linux
	docker build --no-cache . -t tanjunchen-scheduler:$(TAG)

update:
	go mod download
	go mod tidy
	go mod vendor

clean:
	rm -rf _output/
	rm -f *.log